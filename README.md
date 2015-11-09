# Warehouse Manager

This is a simple wrapper for different cloud storage providers (i.e. Google Drive, Dropbox, Amazon S3, etc...).


[![Build Status](https://travis-ci.org/blau-io/warehouse-manager.svg)](https://travis-ci.org/blau-io/warehouse-manager)
[![Coverage Status](https://coveralls.io/repos/blau-io/warehouse-manager/badge.svg?branch=master&service=github)](https://coveralls.io/github/blau-io/warehouse-manager?branch=master)
[![Go Report Card](http://goreportcard.com/badge/blau-io/warehouse-manager)](http://goreportcard.com/report/blau-io/warehouse-manager)

Inspired by such acronyms as *CRUD* or *BREAD*, this API can best be described as *BRRAAP* (after the distinctive sound of a motorcycle): **B**rowse, **R**ead, **R**emove, **A**dd, **A**uthenticate, **P**ublish.

The goal of this API is to abstract all interaction with cloud storage into a single interface. Ideally, an application shouldn't need to worry about where the user stores his/her data. Similar projects exists for Android or iOS as a client library.

**This API is still under heavy development.**

So far, only Google Drive is supported, but support for other providers will follow.
