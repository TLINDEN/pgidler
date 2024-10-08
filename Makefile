
# Copyright © 2023 Thomas von Dein

# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License
# along with this program. If not, see <http://www.gnu.org/licenses/>.


#
# no need to modify anything below
tool = pgidler
dir  = /home/postgres
pod  = pgdkb-servicedbloadte4c2b6-0

all: buildlocal

buildlocal:
	CGO_LDFLAGS='-static' go build -tags osusergo,netgo -ldflags="-extldflags=-static -s"

install: buildlocal
	kubectl cp $(tool) $(pod):$(dir)/$(tool)

clean:
	rm -rf $(tool) coverage.out

goupdate:
	go get -t -u=patch ./...
