# BPlusTree

Go implementation of B+ tree.

## Insertion

```txt
=== RUN   TestInsertNode
    bplustree_test.go:133: tree after insert key: 3
          (3)

    bplustree_test.go:133: tree after insert key: 1
          (1,3)

    bplustree_test.go:133: tree after insert key: 2
              (3)
          (1,2)    (3)

    bplustree_test.go:133: tree after insert key: 4
               (3)
          (1,2)    (3,4)

    bplustree_test.go:133: tree after insert key: 5
                  (3,5)
          (1,2)    (3,4)    (5)

    bplustree_test.go:133: tree after insert key: 6
                   (3,5)
          (1,2)    (3,4)    (5,6)

    bplustree_test.go:133: tree after insert key: 7
                       (5)
               (3)             (7)
          (1,2)    (3,4)    (5,6)    (7)

    bplustree_test.go:133: tree after insert key: 20
                         (5)
               (3)               (7)
          (1,2)    (3,4)    (5,6)    (7,20)

    bplustree_test.go:133: tree after insert key: 18
                             (5)
               (3)                 (7,20)
          (1,2)    (3,4)    (5,6)    (7,18)    (20)

    bplustree_test.go:133: tree after insert key: 19
                               (5,19)
               (3)               (7)              (20)
          (1,2)    (3,4)    (5,6)    (7,18)    (19)    (20)
```
