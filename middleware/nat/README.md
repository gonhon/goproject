# nat网络穿透

## NAT 类型影响：

锥型 NAT：STUN 可穿透。

对称型 NAT：需 TURN 中继。

## STUN/TURN 服务器选择：

公共 STUN 服务器（如 stun.l.google.com:19302）。

自建 TURN 服务器（推荐 coturn）。

## ICE 优化：

优先尝试 STUN，失败后降级到 TURN。

使用 ice.Agent 自动管理候选地址。

## 自建 TURN 服务器（Coturn）
