// +build vendor

package main

import (
	_ "github.com/golang/protobuf/proto"
	_ "google.golang.org/protobuf/reflect/protoreflect"
	_ "google.golang.org/protobuf/runtime/protoimpl"
)

func main() {}
