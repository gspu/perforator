GO_LIBRARY()

LICENSE(BSD-3-Clause)

VERSION(v0.25.0)

SRCS(
    dict.go
    display.go
    lookup.go
    tables.go
)

GO_TEST_SRCS(
    dict_test.go
    display_test.go
)

GO_XTEST_SRCS(examples_test.go)

END()

RECURSE(
    gotest
)
