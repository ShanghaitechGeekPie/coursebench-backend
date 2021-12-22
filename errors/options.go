package errors

type Option interface {
	LogLevel() ErrorLevel
	HideErrorCode() bool
}

type optionImpl struct {
	logLevel      ErrorLevel
	hideErrorCode bool
}

func (o optionImpl) LogLevel() ErrorLevel {
	return o.logLevel
}

func (o optionImpl) HideErrorCode() bool {
	return o.hideErrorCode
}

type OptionBuilder struct {
	opt optionImpl
}

func NewOptionBuilder() *OptionBuilder {
	return &OptionBuilder{
		opt: optionImpl{},
	}
}

func (ob *OptionBuilder) SetLogLevel(level ErrorLevel) *OptionBuilder {
	ob.opt.logLevel = level
	return ob
}

func (ob *OptionBuilder) SetHideErrorCode() *OptionBuilder {
	ob.opt.hideErrorCode = true
	return ob
}

func (ob *OptionBuilder) Build() Option {
	return ob.opt
}
