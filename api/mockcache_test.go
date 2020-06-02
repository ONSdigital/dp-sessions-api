// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package api

import (
	"github.com/ONSdigital/dp-sessions-api/session"
	"sync"
)

var (
	lockCacheMockDeleteAll sync.RWMutex
	lockCacheMockGetByID   sync.RWMutex
	lockCacheMockSet       sync.RWMutex
)

// Ensure, that CacheMock does implement Cache.
// If this is not the case, regenerate this file with moq.
var _ Cache = &CacheMock{}

// CacheMock is a mock implementation of Cache.
//
//     func TestSomethingThatUsesCache(t *testing.T) {
//
//         // make and configure a mocked Cache
//         mockedCache := &CacheMock{
//             DeleteAllFunc: func() error {
// 	               panic("mock out the DeleteAll method")
//             },
//             GetByIDFunc: func(ID string) (*session.Session, error) {
// 	               panic("mock out the GetByID method")
//             },
//             SetFunc: func(s *session.Session) error {
// 	               panic("mock out the Set method")
//             },
//         }
//
//         // use mockedCache in code that requires Cache
//         // and then make assertions.
//
//     }
type CacheMock struct {
	// DeleteAllFunc mocks the DeleteAll method.
	DeleteAllFunc func() error

	// GetByIDFunc mocks the GetByID method.
	GetByIDFunc func(ID string) (*session.Session, error)

	// SetFunc mocks the Set method.
	SetFunc func(s *session.Session) error

	// calls tracks calls to the methods.
	calls struct {
		// DeleteAll holds details about calls to the DeleteAll method.
		DeleteAll []struct {
		}
		// GetByID holds details about calls to the GetByID method.
		GetByID []struct {
			// ID is the ID argument value.
			ID string
		}
		// Set holds details about calls to the Set method.
		Set []struct {
			// S is the s argument value.
			S *session.Session
		}
	}
}

// DeleteAll calls DeleteAllFunc.
func (mock *CacheMock) DeleteAll() error {
	if mock.DeleteAllFunc == nil {
		panic("CacheMock.DeleteAllFunc: method is nil but Cache.DeleteAll was just called")
	}
	callInfo := struct {
	}{}
	lockCacheMockDeleteAll.Lock()
	mock.calls.DeleteAll = append(mock.calls.DeleteAll, callInfo)
	lockCacheMockDeleteAll.Unlock()
	return mock.DeleteAllFunc()
}

// DeleteAllCalls gets all the calls that were made to DeleteAll.
// Check the length with:
//     len(mockedCache.DeleteAllCalls())
func (mock *CacheMock) DeleteAllCalls() []struct {
} {
	var calls []struct {
	}
	lockCacheMockDeleteAll.RLock()
	calls = mock.calls.DeleteAll
	lockCacheMockDeleteAll.RUnlock()
	return calls
}

// GetByID calls GetByIDFunc.
func (mock *CacheMock) GetByID(ID string) (*session.Session, error) {
	if mock.GetByIDFunc == nil {
		panic("CacheMock.GetByIDFunc: method is nil but Cache.GetByID was just called")
	}
	callInfo := struct {
		ID string
	}{
		ID: ID,
	}
	lockCacheMockGetByID.Lock()
	mock.calls.GetByID = append(mock.calls.GetByID, callInfo)
	lockCacheMockGetByID.Unlock()
	return mock.GetByIDFunc(ID)
}

// GetByIDCalls gets all the calls that were made to GetByID.
// Check the length with:
//     len(mockedCache.GetByIDCalls())
func (mock *CacheMock) GetByIDCalls() []struct {
	ID string
} {
	var calls []struct {
		ID string
	}
	lockCacheMockGetByID.RLock()
	calls = mock.calls.GetByID
	lockCacheMockGetByID.RUnlock()
	return calls
}

// Set calls SetFunc.
func (mock *CacheMock) Set(s *session.Session) error {
	if mock.SetFunc == nil {
		panic("CacheMock.SetFunc: method is nil but Cache.Set was just called")
	}
	callInfo := struct {
		S *session.Session
	}{
		S: s,
	}
	lockCacheMockSet.Lock()
	mock.calls.Set = append(mock.calls.Set, callInfo)
	lockCacheMockSet.Unlock()
	return mock.SetFunc(s)
}

// SetCalls gets all the calls that were made to Set.
// Check the length with:
//     len(mockedCache.SetCalls())
func (mock *CacheMock) SetCalls() []struct {
	S *session.Session
} {
	var calls []struct {
		S *session.Session
	}
	lockCacheMockSet.RLock()
	calls = mock.calls.Set
	lockCacheMockSet.RUnlock()
	return calls
}
