/*
Copyright 2021 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package groups

import (
	"strings"

	"gitlab.com/gitlab-org/api/client-go"

	"github.com/crossplane-contrib/provider-gitlab/apis/groups/v1alpha1"
	"github.com/crossplane-contrib/provider-gitlab/pkg/clients"
)

const (
	errMemberNotFound = "404 Group Member Not Found"
)

// MemberClient defines Gitlab Member service operations
type MemberClient interface {
	GetGroupMember(gid interface{}, user int, options ...gitlab.RequestOptionFunc) (*gitlab.GroupMember, *gitlab.Response, error)
	AddGroupMember(gid interface{}, opt *gitlab.AddGroupMemberOptions, options ...gitlab.RequestOptionFunc) (*gitlab.GroupMember, *gitlab.Response, error)
	EditGroupMember(gid interface{}, user int, opt *gitlab.EditGroupMemberOptions, options ...gitlab.RequestOptionFunc) (*gitlab.GroupMember, *gitlab.Response, error)
	RemoveGroupMember(gid interface{}, user int, opt *gitlab.RemoveGroupMemberOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error)
}

// NewMemberClient returns a new Gitlab Group Member service
func NewMemberClient(cfg clients.Config) MemberClient {
	git := clients.NewClient(cfg)
	return git.GroupMembers
}

// IsErrorMemberNotFound helper function to test for errMemberNotFound error.
func IsErrorMemberNotFound(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), errMemberNotFound)
}

// GenerateMemberObservation is used to produce v1alpha1.MemberObservation from
// gitlab.Member.
func GenerateMemberObservation(groupMember *gitlab.GroupMember) v1alpha1.MemberObservation {
	if groupMember == nil {
		return v1alpha1.MemberObservation{}
	}

	o := v1alpha1.MemberObservation{
		Username:          groupMember.Username,
		Name:              groupMember.Name,
		State:             groupMember.State,
		AvatarURL:         groupMember.AvatarURL,
		WebURL:            groupMember.WebURL,
		GroupSAMLIdentity: groupMemberSAMLIdentityGitlabToV1alpha1(groupMember.GroupSAMLIdentity),
	}

	return o
}

// GenerateAddMemberOptions generates group member add options
func GenerateAddMemberOptions(p *v1alpha1.MemberParameters) *gitlab.AddGroupMemberOptions {
	groupMember := &gitlab.AddGroupMemberOptions{
		UserID:      p.UserID,
		AccessLevel: accessLevelValueV1alpha1ToGitlab(&p.AccessLevel),
	}
	if p.ExpiresAt != nil {
		groupMember.ExpiresAt = p.ExpiresAt
	}
	return groupMember
}

// GenerateEditMemberOptions generates group member edit options
func GenerateEditMemberOptions(p *v1alpha1.MemberParameters) *gitlab.EditGroupMemberOptions {
	groupMember := &gitlab.EditGroupMemberOptions{
		AccessLevel: accessLevelValueV1alpha1ToGitlab(&p.AccessLevel),
	}
	if p.ExpiresAt != nil {
		groupMember.ExpiresAt = p.ExpiresAt
	}
	return groupMember
}

// accessLevelValueV1alpha1ToGitlab converts *v1alpha1.AccessLevelValue to *gitlab.AccessLevelValue
func accessLevelValueV1alpha1ToGitlab(from *v1alpha1.AccessLevelValue) *gitlab.AccessLevelValue {
	return (*gitlab.AccessLevelValue)(from)
}

// groupMemberSAMLIdentityGitlabToV1alpha1 converts *gitlab.MemberSAMLIdentity to *v1alpha1.MemberSAMLIdentity
func groupMemberSAMLIdentityGitlabToV1alpha1(from *gitlab.GroupMemberSAMLIdentity) *v1alpha1.MemberSAMLIdentity {
	return (*v1alpha1.MemberSAMLIdentity)(from)
}
