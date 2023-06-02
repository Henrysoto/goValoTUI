package main

// assets/default.png
const artDefault64 = "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"