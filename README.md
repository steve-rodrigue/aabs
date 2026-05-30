# AABS

AABS (Anti-AI Bot Spam) is an open-source browser extension and platform designed to help users identify spam, bot networks, coordinated manipulation, and low-quality engagement on social media.

Rather than focusing solely on whether an account is a bot, AABS analyzes content similarity, posting behavior, semantic patterns, and network activity to detect coordinated campaigns, repetitive social-proof behavior, and other forms of artificial amplification.

AABS assigns trust scores to accounts, posts, comments, and conversations, helping users distinguish authentic discussion from coordinated or deceptive activity.

## Development

### Build and start all services

docker compose up --build 

### Build and start all services in background

docker compose up -d --build 

### View logs

docker compose logs -f 

### Stop all services

docker compose down 

## License

MIT