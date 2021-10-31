package main

check(v) {
	v == "world"
}

response(v, status) = message {
	message = {
		"id": "CUSTOM-123",
		"name": "Does hello equal world",
		"severity": "critical",
		"allow_response": sprintf("'hello' correctly equals '%s'", [v]),
		"deny_response": sprintf("'hello' contains '%s' instead of expected 'world'", [v]),
		"value": v,
		"key": "hello",
		"status": status,
	}
}

deny[msg] {
	v := input[_].hello
	not check(v)
	msg := response(v, "deny")
}

allow[msg] {
	v := input[_].hello
	check(v)
	msg := response(v, "allow")
}
