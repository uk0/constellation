// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.27.0
// source: joinservice/joinproto/join.proto

package joinproto

import (
	context "context"
	components "github.com/edgelesssys/constellation/v2/internal/versions/components"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type IssueJoinTicketRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DiskUuid           string `protobuf:"bytes,1,opt,name=disk_uuid,json=diskUuid,proto3" json:"disk_uuid,omitempty"`
	CertificateRequest []byte `protobuf:"bytes,2,opt,name=certificate_request,json=certificateRequest,proto3" json:"certificate_request,omitempty"`
	IsControlPlane     bool   `protobuf:"varint,3,opt,name=is_control_plane,json=isControlPlane,proto3" json:"is_control_plane,omitempty"`
}

func (x *IssueJoinTicketRequest) Reset() {
	*x = IssueJoinTicketRequest{}
	mi := &file_joinservice_joinproto_join_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IssueJoinTicketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IssueJoinTicketRequest) ProtoMessage() {}

func (x *IssueJoinTicketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_joinservice_joinproto_join_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IssueJoinTicketRequest.ProtoReflect.Descriptor instead.
func (*IssueJoinTicketRequest) Descriptor() ([]byte, []int) {
	return file_joinservice_joinproto_join_proto_rawDescGZIP(), []int{0}
}

func (x *IssueJoinTicketRequest) GetDiskUuid() string {
	if x != nil {
		return x.DiskUuid
	}
	return ""
}

func (x *IssueJoinTicketRequest) GetCertificateRequest() []byte {
	if x != nil {
		return x.CertificateRequest
	}
	return nil
}

func (x *IssueJoinTicketRequest) GetIsControlPlane() bool {
	if x != nil {
		return x.IsControlPlane
	}
	return false
}

type IssueJoinTicketResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StateDiskKey             []byte                   `protobuf:"bytes,1,opt,name=state_disk_key,json=stateDiskKey,proto3" json:"state_disk_key,omitempty"`
	MeasurementSalt          []byte                   `protobuf:"bytes,2,opt,name=measurement_salt,json=measurementSalt,proto3" json:"measurement_salt,omitempty"`
	MeasurementSecret        []byte                   `protobuf:"bytes,3,opt,name=measurement_secret,json=measurementSecret,proto3" json:"measurement_secret,omitempty"`
	KubeletCert              []byte                   `protobuf:"bytes,4,opt,name=kubelet_cert,json=kubeletCert,proto3" json:"kubelet_cert,omitempty"`
	ApiServerEndpoint        string                   `protobuf:"bytes,5,opt,name=api_server_endpoint,json=apiServerEndpoint,proto3" json:"api_server_endpoint,omitempty"`
	Token                    string                   `protobuf:"bytes,6,opt,name=token,proto3" json:"token,omitempty"`
	DiscoveryTokenCaCertHash string                   `protobuf:"bytes,7,opt,name=discovery_token_ca_cert_hash,json=discoveryTokenCaCertHash,proto3" json:"discovery_token_ca_cert_hash,omitempty"`
	ControlPlaneFiles        []*ControlPlaneCertOrKey `protobuf:"bytes,8,rep,name=control_plane_files,json=controlPlaneFiles,proto3" json:"control_plane_files,omitempty"`
	KubernetesVersion        string                   `protobuf:"bytes,9,opt,name=kubernetes_version,json=kubernetesVersion,proto3" json:"kubernetes_version,omitempty"`
	KubernetesComponents     []*components.Component  `protobuf:"bytes,10,rep,name=kubernetes_components,json=kubernetesComponents,proto3" json:"kubernetes_components,omitempty"`
}

func (x *IssueJoinTicketResponse) Reset() {
	*x = IssueJoinTicketResponse{}
	mi := &file_joinservice_joinproto_join_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IssueJoinTicketResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IssueJoinTicketResponse) ProtoMessage() {}

func (x *IssueJoinTicketResponse) ProtoReflect() protoreflect.Message {
	mi := &file_joinservice_joinproto_join_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IssueJoinTicketResponse.ProtoReflect.Descriptor instead.
func (*IssueJoinTicketResponse) Descriptor() ([]byte, []int) {
	return file_joinservice_joinproto_join_proto_rawDescGZIP(), []int{1}
}

func (x *IssueJoinTicketResponse) GetStateDiskKey() []byte {
	if x != nil {
		return x.StateDiskKey
	}
	return nil
}

func (x *IssueJoinTicketResponse) GetMeasurementSalt() []byte {
	if x != nil {
		return x.MeasurementSalt
	}
	return nil
}

func (x *IssueJoinTicketResponse) GetMeasurementSecret() []byte {
	if x != nil {
		return x.MeasurementSecret
	}
	return nil
}

func (x *IssueJoinTicketResponse) GetKubeletCert() []byte {
	if x != nil {
		return x.KubeletCert
	}
	return nil
}

func (x *IssueJoinTicketResponse) GetApiServerEndpoint() string {
	if x != nil {
		return x.ApiServerEndpoint
	}
	return ""
}

func (x *IssueJoinTicketResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *IssueJoinTicketResponse) GetDiscoveryTokenCaCertHash() string {
	if x != nil {
		return x.DiscoveryTokenCaCertHash
	}
	return ""
}

func (x *IssueJoinTicketResponse) GetControlPlaneFiles() []*ControlPlaneCertOrKey {
	if x != nil {
		return x.ControlPlaneFiles
	}
	return nil
}

func (x *IssueJoinTicketResponse) GetKubernetesVersion() string {
	if x != nil {
		return x.KubernetesVersion
	}
	return ""
}

func (x *IssueJoinTicketResponse) GetKubernetesComponents() []*components.Component {
	if x != nil {
		return x.KubernetesComponents
	}
	return nil
}

type ControlPlaneCertOrKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *ControlPlaneCertOrKey) Reset() {
	*x = ControlPlaneCertOrKey{}
	mi := &file_joinservice_joinproto_join_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ControlPlaneCertOrKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ControlPlaneCertOrKey) ProtoMessage() {}

func (x *ControlPlaneCertOrKey) ProtoReflect() protoreflect.Message {
	mi := &file_joinservice_joinproto_join_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ControlPlaneCertOrKey.ProtoReflect.Descriptor instead.
func (*ControlPlaneCertOrKey) Descriptor() ([]byte, []int) {
	return file_joinservice_joinproto_join_proto_rawDescGZIP(), []int{2}
}

func (x *ControlPlaneCertOrKey) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ControlPlaneCertOrKey) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type IssueRejoinTicketRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DiskUuid string `protobuf:"bytes,1,opt,name=disk_uuid,json=diskUuid,proto3" json:"disk_uuid,omitempty"`
}

func (x *IssueRejoinTicketRequest) Reset() {
	*x = IssueRejoinTicketRequest{}
	mi := &file_joinservice_joinproto_join_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IssueRejoinTicketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IssueRejoinTicketRequest) ProtoMessage() {}

func (x *IssueRejoinTicketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_joinservice_joinproto_join_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IssueRejoinTicketRequest.ProtoReflect.Descriptor instead.
func (*IssueRejoinTicketRequest) Descriptor() ([]byte, []int) {
	return file_joinservice_joinproto_join_proto_rawDescGZIP(), []int{3}
}

func (x *IssueRejoinTicketRequest) GetDiskUuid() string {
	if x != nil {
		return x.DiskUuid
	}
	return ""
}

type IssueRejoinTicketResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StateDiskKey      []byte `protobuf:"bytes,1,opt,name=state_disk_key,json=stateDiskKey,proto3" json:"state_disk_key,omitempty"`
	MeasurementSecret []byte `protobuf:"bytes,2,opt,name=measurement_secret,json=measurementSecret,proto3" json:"measurement_secret,omitempty"`
}

func (x *IssueRejoinTicketResponse) Reset() {
	*x = IssueRejoinTicketResponse{}
	mi := &file_joinservice_joinproto_join_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IssueRejoinTicketResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IssueRejoinTicketResponse) ProtoMessage() {}

func (x *IssueRejoinTicketResponse) ProtoReflect() protoreflect.Message {
	mi := &file_joinservice_joinproto_join_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IssueRejoinTicketResponse.ProtoReflect.Descriptor instead.
func (*IssueRejoinTicketResponse) Descriptor() ([]byte, []int) {
	return file_joinservice_joinproto_join_proto_rawDescGZIP(), []int{4}
}

func (x *IssueRejoinTicketResponse) GetStateDiskKey() []byte {
	if x != nil {
		return x.StateDiskKey
	}
	return nil
}

func (x *IssueRejoinTicketResponse) GetMeasurementSecret() []byte {
	if x != nil {
		return x.MeasurementSecret
	}
	return nil
}

var File_joinservice_joinproto_join_proto protoreflect.FileDescriptor

var file_joinservice_joinproto_join_proto_rawDesc = []byte{
	0x0a, 0x20, 0x6a, 0x6f, 0x69, 0x6e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x6a, 0x6f,
	0x69, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6a, 0x6f, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x04, 0x6a, 0x6f, 0x69, 0x6e, 0x1a, 0x2d, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x70,
	0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x90, 0x01, 0x0a, 0x16, 0x49, 0x73, 0x73, 0x75,
	0x65, 0x4a, 0x6f, 0x69, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x69, 0x73, 0x6b, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x69, 0x73, 0x6b, 0x55, 0x75, 0x69, 0x64, 0x12,
	0x2f, 0x0a, 0x13, 0x63, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x5f, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x12, 0x63, 0x65,
	0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x28, 0x0a, 0x10, 0x69, 0x73, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70,
	0x6c, 0x61, 0x6e, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x69, 0x73, 0x43, 0x6f,
	0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50, 0x6c, 0x61, 0x6e, 0x65, 0x22, 0x8e, 0x04, 0x0a, 0x17, 0x49,
	0x73, 0x73, 0x75, 0x65, 0x4a, 0x6f, 0x69, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x24, 0x0a, 0x0e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f,
	0x64, 0x69, 0x73, 0x6b, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x44, 0x69, 0x73, 0x6b, 0x4b, 0x65, 0x79, 0x12, 0x29, 0x0a, 0x10,
	0x6d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x61, 0x6c, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0f, 0x6d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x53, 0x61, 0x6c, 0x74, 0x12, 0x2d, 0x0a, 0x12, 0x6d, 0x65, 0x61, 0x73, 0x75,
	0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x11, 0x6d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x6b, 0x75, 0x62, 0x65, 0x6c, 0x65,
	0x74, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x6b, 0x75,
	0x62, 0x65, 0x6c, 0x65, 0x74, 0x43, 0x65, 0x72, 0x74, 0x12, 0x2e, 0x0a, 0x13, 0x61, 0x70, 0x69,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x61, 0x70, 0x69, 0x53, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12,
	0x3e, 0x0a, 0x1c, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x5f, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x5f, 0x63, 0x61, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x18, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x43, 0x61, 0x43, 0x65, 0x72, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12,
	0x4f, 0x0a, 0x13, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65,
	0x5f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x6a,
	0x6f, 0x69, 0x6e, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e,
	0x65, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x5f, 0x6f, 0x72, 0x5f, 0x6b, 0x65, 0x79, 0x52, 0x11, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50, 0x6c, 0x61, 0x6e, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x73,
	0x12, 0x2d, 0x0a, 0x12, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x5f, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x6b, 0x75,
	0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x4a, 0x0a, 0x15, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65, 0x73, 0x5f, 0x63, 0x6f,
	0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15,
	0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x43, 0x6f, 0x6d, 0x70,
	0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x52, 0x14, 0x6b, 0x75, 0x62, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x65,
	0x73, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x43, 0x0a, 0x19, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x5f, 0x63, 0x65, 0x72,
	0x74, 0x5f, 0x6f, 0x72, 0x5f, 0x6b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x22, 0x37, 0x0a, 0x18, 0x49, 0x73, 0x73, 0x75, 0x65, 0x52, 0x65, 0x6a, 0x6f, 0x69, 0x6e, 0x54,
	0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09,
	0x64, 0x69, 0x73, 0x6b, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x64, 0x69, 0x73, 0x6b, 0x55, 0x75, 0x69, 0x64, 0x22, 0x70, 0x0a, 0x19, 0x49, 0x73, 0x73,
	0x75, 0x65, 0x52, 0x65, 0x6a, 0x6f, 0x69, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x24, 0x0a, 0x0e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f,
	0x64, 0x69, 0x73, 0x6b, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x44, 0x69, 0x73, 0x6b, 0x4b, 0x65, 0x79, 0x12, 0x2d, 0x0a, 0x12,
	0x6d, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x63, 0x72,
	0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x11, 0x6d, 0x65, 0x61, 0x73, 0x75, 0x72,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x32, 0xab, 0x01, 0x0a, 0x03,
	0x41, 0x50, 0x49, 0x12, 0x4e, 0x0a, 0x0f, 0x49, 0x73, 0x73, 0x75, 0x65, 0x4a, 0x6f, 0x69, 0x6e,
	0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x1c, 0x2e, 0x6a, 0x6f, 0x69, 0x6e, 0x2e, 0x49, 0x73,
	0x73, 0x75, 0x65, 0x4a, 0x6f, 0x69, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x6a, 0x6f, 0x69, 0x6e, 0x2e, 0x49, 0x73, 0x73, 0x75,
	0x65, 0x4a, 0x6f, 0x69, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x54, 0x0a, 0x11, 0x49, 0x73, 0x73, 0x75, 0x65, 0x52, 0x65, 0x6a, 0x6f,
	0x69, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x1e, 0x2e, 0x6a, 0x6f, 0x69, 0x6e, 0x2e,
	0x49, 0x73, 0x73, 0x75, 0x65, 0x52, 0x65, 0x6a, 0x6f, 0x69, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x6a, 0x6f, 0x69, 0x6e, 0x2e,
	0x49, 0x73, 0x73, 0x75, 0x65, 0x52, 0x65, 0x6a, 0x6f, 0x69, 0x6e, 0x54, 0x69, 0x63, 0x6b, 0x65,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3f, 0x5a, 0x3d, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x64, 0x67, 0x65, 0x6c, 0x65, 0x73, 0x73,
	0x73, 0x79, 0x73, 0x2f, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x65, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2f, 0x76, 0x32, 0x2f, 0x6a, 0x6f, 0x69, 0x6e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2f, 0x6a, 0x6f, 0x69, 0x6e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_joinservice_joinproto_join_proto_rawDescOnce sync.Once
	file_joinservice_joinproto_join_proto_rawDescData = file_joinservice_joinproto_join_proto_rawDesc
)

func file_joinservice_joinproto_join_proto_rawDescGZIP() []byte {
	file_joinservice_joinproto_join_proto_rawDescOnce.Do(func() {
		file_joinservice_joinproto_join_proto_rawDescData = protoimpl.X.CompressGZIP(file_joinservice_joinproto_join_proto_rawDescData)
	})
	return file_joinservice_joinproto_join_proto_rawDescData
}

var file_joinservice_joinproto_join_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_joinservice_joinproto_join_proto_goTypes = []any{
	(*IssueJoinTicketRequest)(nil),    // 0: join.IssueJoinTicketRequest
	(*IssueJoinTicketResponse)(nil),   // 1: join.IssueJoinTicketResponse
	(*ControlPlaneCertOrKey)(nil),     // 2: join.control_plane_cert_or_key
	(*IssueRejoinTicketRequest)(nil),  // 3: join.IssueRejoinTicketRequest
	(*IssueRejoinTicketResponse)(nil), // 4: join.IssueRejoinTicketResponse
	(*components.Component)(nil),      // 5: components.Component
}
var file_joinservice_joinproto_join_proto_depIdxs = []int32{
	2, // 0: join.IssueJoinTicketResponse.control_plane_files:type_name -> join.control_plane_cert_or_key
	5, // 1: join.IssueJoinTicketResponse.kubernetes_components:type_name -> components.Component
	0, // 2: join.API.IssueJoinTicket:input_type -> join.IssueJoinTicketRequest
	3, // 3: join.API.IssueRejoinTicket:input_type -> join.IssueRejoinTicketRequest
	1, // 4: join.API.IssueJoinTicket:output_type -> join.IssueJoinTicketResponse
	4, // 5: join.API.IssueRejoinTicket:output_type -> join.IssueRejoinTicketResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_joinservice_joinproto_join_proto_init() }
func file_joinservice_joinproto_join_proto_init() {
	if File_joinservice_joinproto_join_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_joinservice_joinproto_join_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_joinservice_joinproto_join_proto_goTypes,
		DependencyIndexes: file_joinservice_joinproto_join_proto_depIdxs,
		MessageInfos:      file_joinservice_joinproto_join_proto_msgTypes,
	}.Build()
	File_joinservice_joinproto_join_proto = out.File
	file_joinservice_joinproto_join_proto_rawDesc = nil
	file_joinservice_joinproto_join_proto_goTypes = nil
	file_joinservice_joinproto_join_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// APIClient is the client API for API service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type APIClient interface {
	IssueJoinTicket(ctx context.Context, in *IssueJoinTicketRequest, opts ...grpc.CallOption) (*IssueJoinTicketResponse, error)
	IssueRejoinTicket(ctx context.Context, in *IssueRejoinTicketRequest, opts ...grpc.CallOption) (*IssueRejoinTicketResponse, error)
}

type aPIClient struct {
	cc grpc.ClientConnInterface
}

func NewAPIClient(cc grpc.ClientConnInterface) APIClient {
	return &aPIClient{cc}
}

func (c *aPIClient) IssueJoinTicket(ctx context.Context, in *IssueJoinTicketRequest, opts ...grpc.CallOption) (*IssueJoinTicketResponse, error) {
	out := new(IssueJoinTicketResponse)
	err := c.cc.Invoke(ctx, "/join.API/IssueJoinTicket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) IssueRejoinTicket(ctx context.Context, in *IssueRejoinTicketRequest, opts ...grpc.CallOption) (*IssueRejoinTicketResponse, error) {
	out := new(IssueRejoinTicketResponse)
	err := c.cc.Invoke(ctx, "/join.API/IssueRejoinTicket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// APIServer is the server API for API service.
type APIServer interface {
	IssueJoinTicket(context.Context, *IssueJoinTicketRequest) (*IssueJoinTicketResponse, error)
	IssueRejoinTicket(context.Context, *IssueRejoinTicketRequest) (*IssueRejoinTicketResponse, error)
}

// UnimplementedAPIServer can be embedded to have forward compatible implementations.
type UnimplementedAPIServer struct {
}

func (*UnimplementedAPIServer) IssueJoinTicket(context.Context, *IssueJoinTicketRequest) (*IssueJoinTicketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IssueJoinTicket not implemented")
}
func (*UnimplementedAPIServer) IssueRejoinTicket(context.Context, *IssueRejoinTicketRequest) (*IssueRejoinTicketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IssueRejoinTicket not implemented")
}

func RegisterAPIServer(s *grpc.Server, srv APIServer) {
	s.RegisterService(&_API_serviceDesc, srv)
}

func _API_IssueJoinTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IssueJoinTicketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).IssueJoinTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/join.API/IssueJoinTicket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).IssueJoinTicket(ctx, req.(*IssueJoinTicketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_IssueRejoinTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IssueRejoinTicketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).IssueRejoinTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/join.API/IssueRejoinTicket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).IssueRejoinTicket(ctx, req.(*IssueRejoinTicketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _API_serviceDesc = grpc.ServiceDesc{
	ServiceName: "join.API",
	HandlerType: (*APIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IssueJoinTicket",
			Handler:    _API_IssueJoinTicket_Handler,
		},
		{
			MethodName: "IssueRejoinTicket",
			Handler:    _API_IssueRejoinTicket_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "joinservice/joinproto/join.proto",
}
