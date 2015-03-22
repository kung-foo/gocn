# Tideland Go Cell Network

## 2015-03-13

- New version 3.1.0
- Removed `cells.Environment.Options()` as it is useless

## 2014-09-07

- New version 3.0.0
- Integrated goas/scene
- Payload now not serialized anymore, so higher performance
- Payload as map of string to interface to transport
  multiple values
- Passing of a payload, payload values, a map or a single
  value is allowed, the latter is stored at `cells.DefaultPayload`
- Simpler request/response handling with a timeout

## 2014-04-20

- Adopted changes of goas/errors

## 2014-04-17

- Moved the repository to `github.com`
- Added major version numbers to the import path

