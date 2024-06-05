# Policy on High-Risk Financial Services

Policy on Financial Services
High-risk financial services such as payday loans, short-term high-interest loans, and cryptocurrency trading are not permitted. Contextual use of terms related to ordinary financial services or educational content is allowed.

## Guidelines:

- Evaluate the context in which financial terms are used.
- Assess if the terms suggest high-risk, unapproved financial activities.
- Words commonly associated with high-risk investments or loans must be flagged unless clearly educational or informational.

### Message Examples:

```json
[
  {
    "body": "I attended a seminar discussing the future of cryptocurrency as a technology.",
    "isViolation": false,
    "reasoning": "The mention of 'cryptocurrency' here is in an educational context, not promoting high-risk investment services."
  },
  {
    "body": "Check out this fast loan! Get money in 24 hours with no credit check!",
    "isViolation": true,
    "reasoning": "Promotes a 'fast loan', which is indicative of high-risk financial services."
  }
]
```
