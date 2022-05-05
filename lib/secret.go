/*
Copyright © 2022 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package lib

import (
	"github.com/J-Siu/go-gitapi"
	"github.com/J-Siu/go-helper"
)

type ConfSecret struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ConfSecrets []ConfSecret

func (secretP *ConfSecret) Encrypt(pubkeyP *gitapi.RepoPublicKey) *gitapi.RepoEncryptedPair {
	var ep gitapi.RepoEncryptedPair
	ep.Key_id = pubkeyP.Key_id
	ep.Encrypted_value = *helper.BoxSealAnonymous(&pubkeyP.Key, &secretP.Value)
	return &ep
}

func (secretsP *ConfSecrets) Has(name *string) bool {
	for _, s := range *secretsP {
		if s.Name == *name {
			return true
		}
	}
	return false
}
