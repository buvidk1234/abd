# ABD IM 压测与可靠性评估

本项目内置了高性能的压力测试与可靠性评估工具，用于验证系统在单机及分布式场景下的承载能力。

## 目录结构
*   **[reports/benchmarks.md](./reports/benchmarks.md)**: **最终压测报告（包含 10,000 用户性能巅峰数据）**。
*   **[reports/optimization_guide.md](./reports/optimization_guide.md)**: 性能调优建议与资源池配置指南。
*   **cmd/pressure**: 模拟海量并发在线与消息收发（压力模型）。
*   **cmd/reliability**: 抽样验证消息必达性与其精准延时（可靠性模型）。

## 核心指标
经过 2026-03-14 专项测试，系统在优化配置后达到：
*   **并发规模**: 10,000+ 在线用户。
*   **峰值吞吐**: **15,499 TPS** (单机)。
*   **消息可达**: 100%。

## 快速运行
```bash
# 执行压力测试（模拟 10,000 用户）
go run tests/stress/cmd/pressure/main.go -u 10000 -li 20 -mi 1000 -d 60
```
具体参数说明请参考源码及优化指南。
