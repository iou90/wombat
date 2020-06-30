// Copyright 2020 Rogchap. All Rights Reserved.

package pb

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
)

type Source interface {
	Services() []string
	Methods() map[string][]*desc.MethodDescriptor
	GetMethodDesc(srv, name string) *desc.MethodDescriptor
}

type fileSource struct {
	files    []*desc.FileDescriptor
	services []string
	methods  map[string][]*desc.MethodDescriptor
}

func (s *fileSource) Services() []string {
	return s.services
}

func (s *fileSource) Methods() map[string][]*desc.MethodDescriptor {
	return s.methods
}

func (s *fileSource) GetMethodDesc(srv, name string) *desc.MethodDescriptor {
	for _, fd := range s.files {
		sdesc := fd.FindService(srv)
		if sdesc == nil {
			continue
		}
		mdesc := sdesc.FindMethodByName(name)
		if mdesc != nil {
			return mdesc
		}
	}
	return nil
}

func GetSourceFromProtoFiles(importPaths, protoPaths []string) (Source, error) {
	filenames, err := protoparse.ResolveFilenames(importPaths, protoPaths...)
	if err != nil {
		return nil, fmt.Errorf("pb: failed to resolve filenames: %v", err)
	}
	parser := protoparse.Parser{
		ImportPaths:      importPaths,
		InferImportPaths: len(importPaths) == 0,
	}
	fds, err := parser.ParseFiles(filenames...)
	if err != nil {
		return nil, fmt.Errorf("pb: failed to parse proto files: %v", err)
	}

	var services []string
	methods := make(map[string][]*desc.MethodDescriptor)
	for _, fd := range fds {
		for _, srv := range fd.GetServices() {
			srvName := srv.GetFullyQualifiedName()
			services = append(services, srvName)
			var ms []*desc.MethodDescriptor
			for _, m := range srv.GetMethods() {
				ms = append(ms, m)
			}
			methods[srvName] = ms
		}
	}

	fmt.Printf("%+v\n", methods)

	return &fileSource{
		files:    fds,
		services: services,
		methods:  methods,
	}, nil
}
