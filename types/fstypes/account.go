package fstypes

// CompanyData represented in files
type AccountData[T any] struct {
	Accounts []T `json:"accounts" yaml:"accounts"`
	Account  T   `json:"account" yaml:"account"`
}

// Merge two data items
func (d *AccountData[T]) Merge(it any) {
	v := it.(*AccountData[T])
	if len(v.Accounts) > 0 {
		d.Accounts = append(d.Accounts, v.Accounts...)
	} else {
		d.Accounts = append(d.Accounts, v.Account)
	}
}

// Result of data as a list
func (d *AccountData[T]) Result() []T {
	return d.Accounts
}

// Reset stored data
func (d *AccountData[T]) Reset() {
	d.Accounts = d.Accounts[:0]
}
