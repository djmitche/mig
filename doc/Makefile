# This Source Code Form is subject to the terms of the Mozilla Public
# # License, v. 2.0. If a copy of the MPL was not distributed with this
# # file, You can obtain one at http://mozilla.org/MPL/2.0/.

RST2HTML	:= rst2html5

doc:
	for doc in $$(ls *.rst); do \
		$(RST2HTML) --stylesheet=docstyle.css "$$doc" > "$$doc.html"; \
	done
	for modname in $(ls ../src/mig/modules/); do \
		if [ -r "../src/mig/modules/${modname}/doc.rst" ]; then \
			$(RST2HTML) --stylesheet=docstyle.css "../src/mig/modules/${modname}/doc.rst" > "module_${modname}.html"; \
		fi; \
	done
	dot -Tsvg -o .files/action_command_flow.svg .files/action_command_flow.dot
