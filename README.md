# README

Simple interface to MongoDB written in Go.

## Key features
1. Simple interface

## Methods
`func (s SimMongoDB) Save(ctx context.Context, c string, data interface{}) (string, error)`

`func (s SimMongoDB) SaveMultiple(ctx context.Context, items map[string]interface{}) ([]string, error)`
SaveMultiple func: c stands for MongoDB collection where data would be saved. e.g save data in 'users' collection in MongoDB
ctx can be a mongodb session context for transactions

`func (s SimMongoDB) GetItem(ctx context.Context, c string, filter map[string]interface{}, excludedFields map[string]interface{}, result interface{}) error`
GetItem func: c stands for collection where item should be retrieved. e.g retrieve item from 'users' collection in MongoDB.
ctx can be a mongodb session context for transactions
results is a pointer to object to store returned data. nil is returned for error if item is found

`func (s SimMongoDB) UpdateItem(ctx context.Context, c string, match map[string]interface{}, update map[string]interface{}) (int64, error)`

`func (s SimMongoDB) UpdateItems(ctx context.Context, c string, match map[string]interface{}, update map[string]interface{}) (int64, error)`

`func (s SimMongoDB) GetItems(ctx context.Context, c string, filter map[string]interface{}, limit int64, excludedFields map[string]interface{}, sort map[string]interface{}, results interface{}) error`
GetItems func: c stands for collection where data would be saved. e.g save data in 'users' collection in MongoDB. id is string
ctx can be a mongodb session context for transactions
results is a pointer to slice of object to store returned data. nil is returned for error if item is found

`func (s SimMongoDB) CountItems(ctx context.Context, c string, filter map[string]interface{}) (int64, error)`
CountItems func: c stands for collection where items should be counted. e.g count items in 'users' collection in MongoDB.
ctx can be a mongodb session context for transactions

`func (s SimMongoDB) DeleteItem(ctx context.Context, c string, filter map[string]interface{}) (int64, error)`
DeleteItem func: c stands for collection where item should be retrieved. e.g retrieve item from 'users' collection in MongoDB.
ctx can be a mongodb session context for transactions


`func (s SimMongoDB) DeleteItems(ctx context.Context, c string, filter map[string]interface{}) (int64, error)`
DeleteItems func: c stands for collection where item should be retrieved. e.g retrieve item from 'users' collection in MongoDB.
ctx can be a mongodb session context for transactions