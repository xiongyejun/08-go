# golang 数据类型

标签（空格分隔）： 数据类型

[原文][1]
---

## basic types
    
    typedef signed char             int8;
    typedef unsigned char           uint8;
    typedef signed short            int16;
    typedef unsigned short          uint16;
    typedef signed int              int32;
    typedef unsigned int            uint32;
    typedef signed long long int    int64;
    typedef unsigned long long int  uint64;
    typedef float                   float32;
    typedef double                  float64;
    
    #ifdef _64BIT
    
    typedef uint64          uintptr;
    typedef int64           intptr;
    typedef int64           intgo; // Go's int
    typedef uint64          uintgo; // Go's uint
    #else
   
    typedef uint32          uintptr;
    typedef int32           intptr;
    typedef int32           intgo; // Go's int
    typedef uint32          uintgo; // Go's uint
    
    #endif
    
## defined types
    
    typedef uint8           bool;
    typedef uint8           byte;

## String
    
    struct String
    {
        byte*   str;
        intgo   len;
    };

## rune
    rune是int32的别名，用于表示unicode字符。


## Slice
    struct  Slice
    {               // must not move anything
        byte*   array;      // actual data
        uintgo  len;        // number of elements
        uintgo  cap;        // allocated number of elements
    };

## interface
接口在golang中的实现比较复杂，在$GOROOT/src/pkg/runtime/type.h中定义了：
   
    struct Type
    {
        uintptr size;
        uint32 hash;
        uint8 _unused;
        uint8 align;
        uint8 fieldAlign;
        uint8 kind;
        Alg *alg;
        void *gc;
        String *string;
        UncommonType *x;
        Type *ptrto;
    };
    
在$GOROOT/src/pkg/runtime/runtime.h中定义了：

    struct Iface
    {
        Itab*   tab;
        void*   data;
    };
    
    struct Eface
    {
        Type*   type;
        void*   data;
    };
    
    struct  Itab
    {
        InterfaceType*  inter;
        Type*   type;
        Itab*   link;
        int32   bad;
        int32   unused;
        void    (*fun[])(void);
    };
    

interface实际上是一个结构体，包括两个成员，一个是指向数据的指针，一个包含了成员的类型信息。Eface是interface{}底层使用的数据结构。因为interface中保存了类型信息，所以可以实现反射。反射其实就是查找底层数据结构的元数据。完整的实现在：$GOROOT/src/pkg/runtime/iface.c 。

## map

golang的map实现是hashtable，源码在：$GOROOT/src/pkg/runtime/hashmap.c 。

    struct Hmap
    {
        uintgo  count;
        uint32  flags;
        uint32  hash0;
        uint8   B;
        uint8   keysize;
        uint8   valuesize;
        uint16  bucketsize;
    
        byte    *buckets;
        byte    *oldbuckets;
        uintptr nevacuate;
    };
    


  [1]: https://www.cnblogs.com/junneyang/p/6203710.html