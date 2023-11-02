// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	twitter "github.com/equimper/twitter"
	mock "github.com/stretchr/testify/mock"
)

// TweetService is an autogenerated mock type for the TweetService type
type TweetService struct {
	mock.Mock
}

// All provides a mock function with given fields: ctx
func (_m *TweetService) All(ctx context.Context) ([]twitter.Tweet, error) {
	ret := _m.Called(ctx)

	var r0 []twitter.Tweet
	if rf, ok := ret.Get(0).(func(context.Context) []twitter.Tweet); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]twitter.Tweet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: ctx, input
func (_m *TweetService) Create(ctx context.Context, input twitter.CreateTweetInput) (twitter.Tweet, error) {
	ret := _m.Called(ctx, input)

	var r0 twitter.Tweet
	if rf, ok := ret.Get(0).(func(context.Context, twitter.CreateTweetInput) twitter.Tweet); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Get(0).(twitter.Tweet)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, twitter.CreateTweetInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateReply provides a mock function with given fields: ctx, parentID, input
func (_m *TweetService) CreateReply(ctx context.Context, parentID string, input twitter.CreateTweetInput) (twitter.Tweet, error) {
	ret := _m.Called(ctx, parentID, input)

	var r0 twitter.Tweet
	if rf, ok := ret.Get(0).(func(context.Context, string, twitter.CreateTweetInput) twitter.Tweet); ok {
		r0 = rf(ctx, parentID, input)
	} else {
		r0 = ret.Get(0).(twitter.Tweet)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, twitter.CreateTweetInput) error); ok {
		r1 = rf(ctx, parentID, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *TweetService) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *TweetService) GetByID(ctx context.Context, id string) (twitter.Tweet, error) {
	ret := _m.Called(ctx, id)

	var r0 twitter.Tweet
	if rf, ok := ret.Get(0).(func(context.Context, string) twitter.Tweet); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(twitter.Tweet)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}