package main

import "errors"

type BundleErr struct {
	errs []error
}

type bundleErr struct {
	errs []error
}

func (be *BundleErr) Bundle() error {
	if len(be.errs) == 0 {
		return nil
	}

	return &bundleErr{
		errs: be.errs,
	}
}

func (be *BundleErr) AddErr(err error) {
	if err != nil {
		be.errs = append(be.errs, err)
	}
}

func (be *bundleErr) Error() string {
	msg := ""

	for _, e := range be.errs {
		msg += ":" + e.Error()
	}

	return msg
}

func (be *bundleErr) GetErrs() []error {
	return be.errs
}

func (be *bundleErr) IsOne(err error) bool {
	for _, e := range be.errs {
		if errors.Is(e, err) {
			return true
		}
	}

	return false
}
