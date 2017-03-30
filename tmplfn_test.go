/**
 * tmplfn are a collection common of functions for use with Go's template
 * packages and with Markdown processor packages.
 *
 * @author R. S. Doiel, <rsdoiel@gmail.com>
 *
 * Copyright (c) 2016, R. S. Doiel
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * * Redistributions of source code must retain the above copyright notice, this
 *   list of conditions and the following disclaimer.
 *
 * * Redistributions in binary form must reproduce the above copyright notice,
 *   this list of conditions and the following disclaimer in the documentation
 *   and/or other materials provided with the distribution.
 *
 * * Neither the name of findfile nor the names of its
 *   contributors may be used to endorse or promote products derived from
 *   this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
package tmplfn

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"text/template"
)

func TestJoin(t *testing.T) {
	m1 := template.FuncMap{
		"helloworld": func() string {
			return "Hello World!"
		},
	}

	m2 := Join(m1, TimeMap, PageMap)
	for _, key := range []string{"year", "helloworld", "nl2p"} {
		if _, OK := m2[key]; OK != true {
			t.Errorf("Can't find %s in m2", key)
		}
	}
}

func TestCodeBlock(t *testing.T) {
	data := map[string]interface{}{
		"data": `
echo "Hello World!"
`,
	}
	tSrc := `
This is a codeblock below

{{codeblock .data 0 0 "shell"}}
`

	expected := fmt.Sprintf(`
This is a codeblock below

%sshell
    echo "Hello World!"
%s
`, "```", "```")

	fMap := Join(TimeMap, PageMap)
	tmpl, err := AssembleString(fMap, tSrc)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	buf := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(buf, data)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	result := fmt.Sprintf("%s", buf)
	if result != expected {
		t.Errorf("codeblock expected:\n\n%s\n\ngot:\n\n%s\n", expected, buf)
		t.FailNow()
	}

	data["data"] = `
# This is a comment.
if [[ i > "1" ]]; then
	echo "i is $i"
fi

# done!
`
	expected = fmt.Sprintf(`
This is a codeblock below

%sshell
    # This is a comment.
    if [[ i > 1 ]]; then
        echo "i is $i"
    fi 

    # done!
%s
`, "```", "```")

	buf = bytes.NewBuffer([]byte{})
	err = tmpl.Execute(buf, data)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}

	result = fmt.Sprintf("%s", buf)
	if strings.Compare(result, expected) == 0 {
		t.Errorf("codeblock expected:\n\n[%s]\n\ngot:\n\n[%s]\n", expected, buf)
		t.FailNow()
	}

}
