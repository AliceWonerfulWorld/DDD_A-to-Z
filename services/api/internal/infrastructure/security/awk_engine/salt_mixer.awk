{
	if (NR > 1) {
		input = input "\n"
	}
	input = input $0
}

END {
	for (i = length(input); i >= 1; i--) {
		reversed = reversed substr(input, i, 1)
	}

	for (i = 1; i <= length(reversed); i++) {
		output = output substr(reversed, i, 1)
		if (i % 3 == 0) {
			output = output salt_char(i)
		}
	}

	printf "%s::%d:%d", output, length(input), length(salt)
}

function salt_char(pos, salt_length, salt_index) {
	salt_length = length(salt)
	if (salt_length == 0) {
		return ""
	}

	salt_index = (pos % salt_length) + 1
	return substr(salt, salt_index, 1)
}
