# 排错记录 (Troubleshooting Log)

> 本文件记录项目开发过程中遇到的问题、根因和解决方案，方便日后回顾、避免重复踩坑。

---

## 1. 测试报错 `transaction has no signature`（并伴随 gob panic）

**日期**：2026-08-20

**问题现象**：
- `TestVerifyBlock`、`TestAddBlock`、`TestGetHeader` 全部失败，报 `transaction has no signature`。
- `TestAddBlock` 随后还 panic：`gob: cannot encode nil pointer of type *core.Header`。

**根因**：
- `Block.Verify()`（`core/block.go`）会遍历区块内每一笔交易并校验签名：
  ```go
  for _, tx := range b.Transactions {
      if err := tx.Verify(); err != nil {   // 未签名的交易返回 "transaction has no signature"
          return err
      }
  }
  ```
- 但测试辅助函数 `randomBlock` / `randomBlockWithPreBlockHash` 构造区块时塞了两笔**没有签名**的占位交易：
  ```go
  txs := []Transaction{
      {Data: []byte("test block1")},
      {Data: []byte("test block2")},
  }
  ```
- 因此所有走 `Verify()` 的区块都会因为占位交易未签名而失败。
- panic 是连锁反应：第一块没加进去 → 链高度停在 0 → `getPrevBlockHash` 调 `GetHeader(1)` 拿到 nil → 对 nil `*Header` 做 gob 编码 panic。

**解决**：
- 去掉 `randomBlock` / `randomBlockWithPreBlockHash` 里无签名的占位交易，改为 `NewBlock(header, nil)`。
- `TestVerifyBlock` 里把 `b.Height = 200` 改掉后，最后恢复时要同时恢复 `b.Height = 0`（只恢复 `Validator` 不够，header 变了签名就不匹配）。

---

## 2. gob 序列化报错 `elliptic.nistCurve[...] has no exported fields`

**日期**：2026-08-20

**问题现象**：
- `TestTxEncodeDecode` 失败，报：
  ```
  gob: type elliptic.nistCurve[*crypto/internal/fips140/nistec.P256Point] has no exported fields
  ```

**根因**：
- `Transaction.From` 字段是 `crypto.PublicKey`，内部是 `*ecdsa.PublicKey`，而 `ecdsa.PublicKey.Curve` 的具体类型 `*elliptic.nistCurve` 没有任何导出字段，gob 无法序列化它的值。
- `gob.Register(elliptic.P256())` 无效：Register 只是给接口字段注册类型名，不能让「没有导出字段」的类型变得可序列化。

**解决**：
- 给 `crypto.PublicKey` 实现 `GobEncode` / `GobDecode`，把公钥序列化成压缩点字节，绕开对 `Curve` 字段的序列化：
  ```go
  func (k PublicKey) GobEncode() ([]byte, error) {
      return elliptic.MarshalCompressed(k.Key.Curve, k.Key.X, k.Key.Y), nil
  }

  func (k *PublicKey) GobDecode(data []byte) error {
      curve := elliptic.P256()
      x, y := elliptic.UnmarshalCompressed(curve, data)
      if x == nil {
          return fmt.Errorf("invalid public key data")
      }
      k.Key = &ecdsa.PublicKey{Curve: curve, X: x, Y: y}
      return nil
  }
  ```
- 删除 `core/encoding.go` 里无效的 `gob.Register(elliptic.P256())` 及 `crypto/elliptic` 导入。

**经验总结**：
- 涉及 `ecdsa.PublicKey` / `ecdsa.PrivateKey` / `elliptic.Curve` 的 gob 序列化，必须给包裹它们的自定义类型实现 `GobEncode`/`GobDecode`（用 `MarshalCompressed`/`UnmarshalCompressed` 或 `Bytes`/`Parse` 系列 API），而不是 `gob.Register`。
