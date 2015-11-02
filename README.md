# Warehouse Manager

This is a simple wrapper for different cloud storage providers (i.e. Google Drive, Dropbox, Amazon S3, etc...).
Inspired by such acronyms as *CRUD* or *BREAD*, this API can best be described as *BRRAAP* (after the distinctive sound of a motorcycle): **B**rowse, **R**ead, **R**emove, **A**uthenticate, **A**dd, **P**ublish.

The goal of this API is to abstract all interaction with cloud storage into a single interface. Ideally, an application shouldn't need to worry about where the user stores his/her data. Similar projects exists for Android or iOS as a client library.

**This API is still under heavy development.**

So far, only Google Drive is supported, but support for other providers will follow.
