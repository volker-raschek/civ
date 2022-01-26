package usecases

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"git.cryptic.systems/volker.raschek/civ/pkg/domain"
	"github.com/Masterminds/semver/v3"
)

// ContainerRuntime is an interface for different container runtimes to return labels
// based on their full qualified container image name. For example:
//
//   imageLabels, err := Load(ctx, "docker.io/library/alpine:latest")
//   imageLabels, err := Load(ctx, "docker.io/library/busybox:latest")
type ContainerRuntime interface {
	GetImageLabels(ctx context.Context, name string) (map[string]string, error)
}

type ConfigLoader interface {
	Load(ctx context.Context) (*domain.Config, error)
}

type LabelVerifier struct {
	config      *domain.Config
	labelLoader ContainerRuntime
	labelStore  *labelStore
}

// Run start the verification process based on the passed config.
func (lv *LabelVerifier) Run(ctx context.Context) error {
	if err := lv.fillLabelStore(ctx); err != nil {
		return err
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(lv.config.Images))

	for image := range lv.config.Images {
		go func(image string) {
			defer func() { wg.Done() }()
			lv.runLabelConstraints(image)
		}(image)
	}

	wg.Wait()

	return nil
}

func (lv *LabelVerifier) runLabelConstraints(image string) {
	for labelKey, labelConstraint := range lv.config.Images[image].LabelConstraints {
		// fetch existing labels from store
		existingLabels := lv.labelStore.GetLabelsForImage(image)

		switch {
		case strings.HasPrefix(labelKey, "%") && strings.HasSuffix(labelKey, "%"):
			m := strings.TrimPrefix(strings.TrimSuffix(labelKey, "%"), "%")

			re, err := regexp.Compile(fmt.Sprintf("^.*%v.*$", m))
			if err != nil {
				labelConstraint.CountResultMessage = err.Error()
			} else {
				state := labelCount(re, labelConstraint.Count, existingLabels)
				labelConstraint.CountResult = &state
			}
		case strings.HasPrefix(labelKey, "%"):
			m := strings.TrimPrefix(labelKey, "%")

			re, err := regexp.Compile(fmt.Sprintf("^.*%v$", m))
			if err != nil {
				labelConstraint.CountResultMessage = err.Error()
			} else {
				state := labelCount(re, labelConstraint.Count, existingLabels)
				labelConstraint.CountResult = &state
			}
		case strings.HasSuffix(labelKey, "%"):
			m := strings.TrimSuffix(labelKey, "%")

			re, err := regexp.Compile(fmt.Sprintf("^%v.*$", m))
			if err != nil {
				labelConstraint.CountResultMessage = err.Error()
			} else {
				state := labelCount(re, labelConstraint.Count, existingLabels)
				labelConstraint.CountResult = &state
			}
		default:
			// labelExists
			if labelConstraint.Exists != nil {
				state := lv.labelExists(labelKey, existingLabels)
				labelConstraint.ExistsResult = state
				if state {
					labelConstraint.ExistsResultMessage = "Label found"
				} else {
					labelConstraint.ExistsResultMessage = "Label not found"
				}
			}

			// labelCompareSemver
			if labelConstraint.CompareSemver != nil {
				labelExistState := lv.labelExists(labelKey, existingLabels)
				if labelExistState {
					parsedSemVer, err := semver.NewVersion(existingLabels[labelKey])
					if err != nil {
						b := false
						labelConstraint.CompareSemverResult = &b
						labelConstraint.CompareSemverResultMessage = err.Error()
					}

					state := labelCompareSemver(labelConstraint.CompareSemver, parsedSemVer)
					labelConstraint.CompareSemverResult = &state
				} else {
					labelConstraint.CompareSemverResult = &labelExistState
					labelConstraint.CompareSemverResultMessage = "Label found"
				}
			}

			// labelCompareString
			if labelConstraint.CompareString != nil {
				state := labelCompareString(labelConstraint.CompareString, existingLabels[labelKey])
				labelConstraint.CompareStringResult = &state
			}
		}
	}
}

func labelCompareSemver(compareSemver *domain.LabelConstraintCompareSemver, parsedSemVer *semver.Version) bool {
	var majorState bool

	// Equal
	if compareSemver.Equal != "" {
		compareSemverEqualVersion, err := semver.NewVersion(compareSemver.Equal)
		if err != nil {
			compareSemver.EqualResultMessage = err.Error()
		} else {
			state := parsedSemVer.Equal(compareSemverEqualVersion)
			compareSemver.EqualResult = &state
			if state {
				compareSemver.EqualResultMessage = fmt.Sprintf("Version %s is equal to %s", parsedSemVer.String(), compareSemverEqualVersion.String())
			} else {
				compareSemver.EqualResultMessage = fmt.Sprintf("Version %s is not equal to %s", parsedSemVer.String(), compareSemverEqualVersion.String())
			}
		}
	}

	// GreaterThan
	if compareSemver.GreaterThan != "" {
		compareSemverGreaterThanVersion, err := semver.NewVersion(compareSemver.GreaterThan)
		if err != nil {
			compareSemver.GreaterThanResultMessage = err.Error()
		} else {
			state := parsedSemVer.GreaterThan(compareSemverGreaterThanVersion)
			compareSemver.GreaterThanResult = &state
			if state {
				compareSemver.GreaterThanResultMessage = fmt.Sprintf("Version %s is greater than %s", parsedSemVer.String(), compareSemverGreaterThanVersion.String())
			} else {
				compareSemver.GreaterThanResultMessage = fmt.Sprintf("Version %s is not greater than %s", parsedSemVer.String(), compareSemverGreaterThanVersion.String())
			}
		}
	}

	// LessThan
	if compareSemver.LessThan != "" {
		compareSemverLessThanVersion, err := semver.NewVersion(compareSemver.LessThan)
		if err != nil {
			compareSemver.LessThanResultMessage = err.Error()
		} else {
			state := parsedSemVer.LessThan(compareSemverLessThanVersion)
			compareSemver.LessThanResult = &state
			if state {
				compareSemver.LessThanResultMessage = fmt.Sprintf("Version %s is lower than %s", parsedSemVer.String(), compareSemverLessThanVersion.String())
			} else {
				compareSemver.LessThanResultMessage = fmt.Sprintf("Version %s is not lower than %s", parsedSemVer.String(), compareSemverLessThanVersion.String())
			}
		}
	}

	return majorState
}

func labelCompareString(compareString *domain.LabelConstraintCompareString, labelValue string) bool {
	var majorState bool = true

	// Equal
	if compareString.Equal != "" {
		state := compareString.Equal == labelValue
		if compareString.Equal == labelValue {
			compareString.EqualResult = &state
			compareString.EqualResultMessage = fmt.Sprintf("%s and %s are equal", labelValue, compareString.Equal)
		} else {
			compareString.EqualResult = &state
			compareString.EqualResultMessage = fmt.Sprintf("%s and %s are not equal", labelValue, compareString.Equal)
		}
	}

	// hasPrefix
	if compareString.HasPrefix != "" {
		state := strings.HasPrefix(labelValue, compareString.HasPrefix)
		if state {
			compareString.HasPrefixResult = &state
			compareString.HasPrefixResultMessage = fmt.Sprintf("%s has prefix %s", labelValue, compareString.HasPrefix)
		} else {
			compareString.HasPrefixResult = &state
			compareString.HasPrefixResultMessage = fmt.Sprintf("%s has not prefix %s", labelValue, compareString.HasPrefix)
		}
	}

	// hasSuffix
	if compareString.HasSuffix != "" {
		state := strings.HasSuffix(labelValue, compareString.HasSuffix)
		if state {
			compareString.HasSuffixResult = &state
			compareString.HasSuffixResultMessage = fmt.Sprintf("%s has suffix %s", labelValue, compareString.HasSuffix)
		} else {
			compareString.HasSuffixResult = &state
			compareString.HasSuffixResultMessage = fmt.Sprintf("%s has not suffix %s", labelValue, compareString.HasSuffix)
		}
	}

	return majorState
}

func labelCount(re *regexp.Regexp, labelConstraintCounter *domain.LabelConstraintCounter, labels map[string]string) bool {
	var majorState bool = true

	var i uint = 0
	for key := range labels {
		if re.MatchString(key) {
			i++
		}
	}

	switch {
	case labelConstraintCounter.Equal != nil:
		switch {
		case i == *labelConstraintCounter.Equal:
			state := true
			labelConstraintCounter.EqualResult = &state
		case i > *labelConstraintCounter.Equal:
			fallthrough
		case i < *labelConstraintCounter.Equal:
			state := false
			labelConstraintCounter.EqualResult = &state
			labelConstraintCounter.EqualResultMessage = fmt.Sprintf("%v is not equal %v", i, *labelConstraintCounter.Equal)
			majorState = false
		}
	case labelConstraintCounter.LessThan != nil:
		switch {
		case i < *labelConstraintCounter.LessThan:
			state := true
			labelConstraintCounter.LessThanResult = &state
		case i >= *labelConstraintCounter.LessThan:
			state := false
			labelConstraintCounter.LessThanResult = &state
			labelConstraintCounter.LessThanResultMessage = fmt.Sprintf("%v is not less than %v", i, *labelConstraintCounter.Equal)
			majorState = false
		}
	case labelConstraintCounter.GreaterThan != nil:
		switch {
		case i < *labelConstraintCounter.GreaterThan:
			state := true
			labelConstraintCounter.GreaterThanResult = &state
		case i >= *labelConstraintCounter.GreaterThan:
			state := false
			labelConstraintCounter.GreaterThanResult = &state
			labelConstraintCounter.GreaterThanResultMessage = fmt.Sprintf("%v is not greater than %v", i, *labelConstraintCounter.Equal)
			majorState = false
		}
	}

	return majorState
}

func (lv *LabelVerifier) labelExists(requiresLabelKey string, existingLabels map[string]string) bool {
	for existingLabelKey := range existingLabels {
		if existingLabelKey == requiresLabelKey {
			return true
		}
	}
	return false
}

// fillLabelStore fills the label store with the labels of the defined images
// from config.
func (lv *LabelVerifier) fillLabelStore(ctx context.Context) error {
	wg := new(sync.WaitGroup)
	wg.Add(len(lv.config.Images))

	errorChannel := make(chan error, len(lv.config.Images))

	for image := range lv.config.Images {
		go func(image string) {
			defer wg.Done()
			labels, err := lv.labelLoader.GetImageLabels(ctx, image)
			if err != nil {
				errorChannel <- err
				return
			}
			lv.labelStore.AddLabelsForImage(image, labels)
		}(image)
	}

	wg.Wait()
	close(errorChannel)

	for {
		err, open := <-errorChannel
		if err != nil {
			return err
		}
		if !open {
			break
		}
	}

	return nil
}

func NewLabelVerifier(config *domain.Config, labelLoader ContainerRuntime) (*LabelVerifier, error) {
	return &LabelVerifier{
		config:      config,
		labelLoader: labelLoader,
		labelStore:  newLabelStore(),
	}, nil
}

type labelStore struct {
	labels map[string]map[string]string
	mutex  *sync.RWMutex
}

func (ls *labelStore) AddLabelsForImage(image string, labels map[string]string) {
	ls.mutex.Lock()
	defer func() { ls.mutex.Unlock() }()
	ls.labels[image] = labels
}

func (ls *labelStore) GetLabelsForImage(image string) map[string]string {
	ls.mutex.RLock()
	defer func() { ls.mutex.RUnlock() }()
	return ls.labels[image]
}

func newLabelStore() *labelStore {
	return &labelStore{
		labels: make(map[string]map[string]string),
		mutex:  new(sync.RWMutex),
	}
}
