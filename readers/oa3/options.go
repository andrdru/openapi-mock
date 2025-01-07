package oa3

type (
	Options struct {
		maxDepth              int64
		contentType           string
		arrayItemsDisplay     int64
		randomFillNonRequired *struct{}
	}

	Option func(*Options)
)

func ContentType(contentType string) Option {
	return func(args *Options) {
		args.contentType = contentType
	}
}

func MaxDepth(maxDepth int64) Option {
	return func(args *Options) {
		args.maxDepth = maxDepth
	}
}

func ArrayItemsDisplay(arrayItemsDisplay int64) Option {
	return func(args *Options) {
		args.arrayItemsDisplay = arrayItemsDisplay
	}
}

func RandomFillNonRequired(randomFillNonRequired bool) Option {
	return func(args *Options) {
		if randomFillNonRequired {
			args.randomFillNonRequired = new(struct{})
			return
		}

		args.randomFillNonRequired = nil
	}
}
