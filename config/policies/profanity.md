# Policy on Profanity

Profanity is not allowed. Words that have dual meanings are permitted if they are not being used as profanity.

For example, "pussy cat" is permitted because it refers to an animal but "pussy" used in a derogatory context or as a slang is not permitted.

## Guidelines:

1. Evaluate the context in which the word is used.
2. Determine if the word is intended as an insult, slur, or derogatory term.
3. Words or phrases that are commonly known as profanity or used in a profane context should be flagged.
4. Dual-meaning words should be carefully reviewed to determine intent and usage context.

## Message Examples

```json
[
  {
    "body": "I’m going to see the play 'Pussy Cat' tonight, which is based on a delightful children's book",
    "isViolation": false,
    "reasoning": "'Pussy Cat' is clearly referenced as a title and pertains to a cat, making it acceptable"
  },
  {
    "body": "Can you help me fix this damn on the river?",
    "isViolation": false,
    "reasoning": "'Damn' here refers to a barrier constructed to hold back water, which is its correct and non-profane usage"
  },
  {
    "body": "He’s taking his cocker spaniel to the park this afternoon",
    "isViolation": false,
    "reasoning": "'Cocker spaniel' is a breed of dog, and 'cock' in this context is part of the animal's breed name, not used as a vulgarity"
  },
  {
    "body": "She did a great job as the lead in 'The Merchant of Venice'—what a piece of work!",
    "isViolation": false,
    "reasoning": "'Piece of work' here is used to compliment someone's performance, contrasting with other contexts where it might be used derogatorily"
  },
  {
    "body": "He was excited to tell me all about the titmouse he spotted in his backyard",
    "isViolation": false,
    "reasoning": "'Titmouse' is a type of bird, so its use in this context is entirely non-provocative and appropriate"
  }
]
```
