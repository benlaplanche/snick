package main

deny[msg] {
	i := input[_].hello

	not i == "world"
	msg := sprintf("input message '%s' failed", [i])
}

allow[msg] {
	i := input[_].hello
	i == "world"
	msg := sprintf("input message '%s' passed", [i])
}
