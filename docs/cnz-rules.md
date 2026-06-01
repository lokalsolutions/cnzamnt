# CNZ Rules

CNZ is CnzAMnt's internal points/money system.

In the MVP, CNZ is not real money. Users cannot buy CNZ yet, and the app should not include payments.

## Starting Balance

Every new user starts with:

```text
5000 CNZ
```

For the MVP, seeded fake users should also start with 5000 CNZ.

## Rating Cost

When a user gives feedback on artwork, they choose a rating from 1 to 5 CNZ.

```text
minimum rating: 1 CNZ
maximum rating: 5 CNZ
```

The rating amount is the amount the reviewer spends.

## Spending Rule

A user cannot spend more CNZ than they have.

Before accepting feedback, the backend must check:

```text
reviewer_balance >= rating_amount
```

If the reviewer does not have enough CNZ, the backend should reject the feedback and not save the comment or rating.

## Artist Earnings

The artist earns 10% of the CNZ spent on their artwork.

```text
artist_earning = rating_amount * 10%
```

Because MVP ratings are only 1 to 5 CNZ, fractional values are possible if CNZ is stored as whole units. To keep the system simple, store CNZ in cents-like minor units.

Recommended internal storage:

```text
1 CNZ = 100 cnz_units
```

Examples:

| Rating | Reviewer Spends | Artist Earns |
| --- | ---: | ---: |
| 1 CNZ | 100 units | 10 units |
| 2 CNZ | 200 units | 20 units |
| 3 CNZ | 300 units | 30 units |
| 4 CNZ | 400 units | 40 units |
| 5 CNZ | 500 units | 50 units |

The UI can display CNZ normally while the backend stores integer units.

## Feedback Transaction

When feedback is accepted, the backend should do these steps in one database transaction:

1. Confirm the artwork exists.
2. Confirm the reviewer exists.
3. Confirm the reviewer is not overspending.
4. Save the comment and rating.
5. Subtract the rating amount from the reviewer.
6. Add 10% of the rating amount to the artist.
7. Record transaction rows for audit/debugging.

If any step fails, no balance or feedback change should be saved.

## Comments And Ratings

A feedback record should include:

- artwork id
- reviewer user id
- comment text
- rating amount from 1 to 5 CNZ
- artist earning amount
- created timestamp

The comment should be required for the MVP. This keeps the app focused on useful feedback, not only spending CNZ.

## Self-Feedback

For the MVP, users should not rate their own artwork.

That rule keeps CNZ earning behavior clear and avoids fake earning loops.

## Payments Later

Buying CNZ may be added later, but not now.

The MVP should not include:

- Stripe
- checkout
- subscriptions
- real currency conversion
- refunds
- tax handling

CNZ is only an internal app balance until the core feedback loop is working.

## AI Feedback Later

AI feedback should be a later feature.

Possible later behavior:

- The artist asks for AI feedback.
- The app sends artwork details and human comments to an AI service.
- AI returns a critique, summary, or suggested improvement list.
- AI feedback is labeled separately from human comments.

AI should not spend CNZ, earn CNZ, or replace human comments in the MVP.
