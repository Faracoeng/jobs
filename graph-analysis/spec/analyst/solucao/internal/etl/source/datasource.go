package source

type DataSource interface {
    ReadAll() ([][]string, error)
}
