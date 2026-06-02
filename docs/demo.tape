Output docs/demo.gif

Set Shell "bash"
Set FontSize 14
Set Width 900
Set Height 550
Set Theme "Catppuccin Mocha"
Set Padding 20
Set PlaybackSpeed 1.25

# Binary and vault are pre-seeded before vhs runs (see Makefile/record.sh).
# ENVII_PASSPHRASE is set in the environment.

# --- Open envii ---
Type "ENVII_PASSPHRASE=demo /tmp/envii-demo --vault /tmp/demo-vault.age"
Enter
Sleep 2s

# --- Select first project (my-api) ---
Enter
Sleep 700ms

# --- Select dev env ---
Enter
Sleep 700ms

# --- Browse vars ---
Down
Sleep 300ms
Down
Sleep 300ms

# --- Reveal secret ---
Type "r"
Sleep 900ms
Type "r"
Sleep 500ms

# --- Copy to clipboard ---
Up
Sleep 300ms
Type "c"
Sleep 800ms

# --- Add a new var ---
Type "a"
Sleep 500ms
Type "DEBUG"
Enter
Sleep 400ms
Type "true"
Enter
Sleep 700ms

# --- Save ---
Type "s"
Sleep 900ms

# --- Back to projects ---
Escape
Sleep 400ms
Escape
Sleep 500ms

# --- Navigate to second project ---
Down
Sleep 400ms
Enter
Sleep 700ms

# --- Quit ---
Type "q"
Sleep 600ms
