# Compute Cartel

A little terminal game I built in Go where you play as the CEO of a cloud computing startup. 

The goal is to survive 12 quarters of pricing wars against an AI rival, deal with random market crashes, and try not to go bankrupt.  

It uses a Game Theory algorithm (Tit-for-Tat). If you play nice, it plays nice. If you undercut it, it *will* take revenge.

You can't actually see how much money the AI has unless you pay ₹3M to spy on them.

When you think you've got the AI figured out, a massive 3rd-party megacorp enters the game on Quarter 4 and starts stealing everyone's customers. 

You can spend cash to upgrade your tech for permanent bonuses, or string together hidden move combos (like *Match -> Match -> Undercut*) to land a massive critical hit on the market.

A scrolling marquee throws random events at you like power grid failures or viral tech drops that completely change the math.

## 🚀 How to play it

Clone the repository and make sure you have Go installed on your system.

```bash 

git clone https://github.com/abhisheksunil2201/compute-cartel.git

cd compute-cartel

go mod tidy

go run main.go

```

You have 12 Quarters to either make more money than the opponent, or drive them into the ground.

Make your move (The system plays against you at the same time):

1. `[u] Undercut` - Slash prices to steal customers (but nobody makes much money).

2. `[m] Match` - Keep prices normal and stable.

3. `[p] Premium` - Charge luxury prices. Great money, if the AI doesn't undercut you.

Spend your cash (Free actions):

4. `[i] Invest (₹10M)` - Upgrade your tech so your moves make more money permanently.

5. `[s] Scout (₹3M)` - Pay some hackers to reveal the AI's hidden stats.

Other stuff:

6. `[q]` - Quit the game.

7. `[r]` - Restart from Quarter 0.
