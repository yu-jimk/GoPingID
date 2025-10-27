# GoPingID
ICMP Identifierを指定できるGo製Pingツール。

## 特徴
- ICMP Echo Requestを送信
- ICMP Identifierを任意で指定可能（PIDとは独立）
- RTTをミリ秒単位で表示

## 使い方
```bash
sudo go run main.go -a 8.8.8.8 -id 2 -n 3
```

## オプション
| オプション | デフォルト   | 説明                        | 
| ---------- | :----------- | :-------------------------- | 
| `-a`       | (必須)       | 送信先IPアドレス            | 
| `-n`       | `3`            | Echo Request送信回数        | 
| `-id`      | `PID & 0xffff` | ICMP Identifier（`0〜65535`） | 
| `-t`       | `3s`           | タイムアウト時間            | 

## 出力例
```bash
PING 8.8.8.8 (id=2):
8.8.8.8: icmp_seq=0 id=2 time=10ms
8.8.8.8: icmp_seq=1 id=2 time=18ms
8.8.8.8: icmp_seq=2 id=2 time=20ms
```
