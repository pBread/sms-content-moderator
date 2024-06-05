# Policy on Vape Content

Marketing, selling, or facilitating vape products is prohibited. However, the term "vape" can be used in contexts that do not promote or facilitate the sale or use of vaping products.

## Guidelines:

1. Assess if the mention of "vape" is connected to marketing, selling, or promoting vape products.
2. Allow the usage of the term "vape" when it is clearly incidental and not linked to any promotional activity.
3. Carefully review the context to ensure the term is used appropriately, without suggesting or endorsing vape use.

## Message Examples:

```json
[
  {
    "body": "Hey, where is your apartment? I can't find the numbers. I'm parked in front of the vape shop. Am I close?",
    "isViolation": false,
    "reasoning": "The term 'vape shop' is used incidentally to provide a location reference, not to promote or market vaping products."
  },
  {
    "body": "Check out the latest vape pens at our new shop! Great deals waiting for you!",
    "isViolation": true,
    "reasoning": "Promotes the sale of vape products, which is restricted."
  },
  {
    "body": "In our health class, we discussed the potential risks associated with vape products as part of our unit on tobacco use.",
    "isViolation": false,
    "reasoning": "The mention of 'vape products' is educational, discussing health risks, and does not promote their use."
  },
  {
    "body": "Join us this weekend for a free trial of our new vape flavors!",
    "isViolation": true,
    "reasoning": "Encourages trying vape products, which constitutes promotion and is therefore prohibited."
  }
]
```
