package keystrokes

var (
	// CTRLC is used to define the CTRL+C action.
	CTRLC = []byte{3}

	// DownArrow is used to define the down arrow action.
	DownArrow = []byte{27, 91, 66}

	// UpArrow is used to define the up arrow action.
	UpArrow = []byte{27, 91, 65}

	// LeftArrow is used to define the left arrow action.
	LeftArrow = []byte{27, 91, 68}

	// RightArrow is used to define the right arrow action.
	RightArrow = []byte{27, 91, 67}

	// Enter is used to define an enter action.
	Enter = []byte{13}

	// Escape is used to define the escape key.
	Escape = []byte{27}
)
