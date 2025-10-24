# Troubleshooting

## Port Already in Use
```bash
# Check port usage
lsof -i :8080  # Backend
lsof -i :3000  # Frontend Docker
lsof -i :5173  # Frontend dev

# Modify ports in docker-compose.yml or vite.config.ts
```

## Go Module Issues
```bash
# Set proxy (China)
export GOPROXY=https://goproxy.cn,direct

# Clean and reinstall
cd backend
rm go.sum
go mod tidy
```

## Docker Build Issues
```bash
# Clean Docker cache
docker system prune -a

# Rebuild without cache
docker-compose build --no-cache
```

## Frontend Connection Issues
Check backend is running: `curl http://localhost:8080/health`
Verify `VITE_API_BASE_URL` in `frontend/.env`
