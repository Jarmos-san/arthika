// Package handler_test contains black-box tests for HTTP handler.
//
// These tests validate the behaviour of the transport layer in isolation by:
//   - mocking service dependencies
//   - issuing HTTP requests using httptest utilities
//   - asserting on HTTP responses (status, headers, body)
//
// The goal is NOT to test business logic, but to ensure correct interaction
// between the handler and its dependencies, along with proper HTTP semantics.
package handler_test
