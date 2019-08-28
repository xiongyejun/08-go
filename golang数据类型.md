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
接口在golang中的实现比较复杂，在$GOROOT/src/runtime/type.h中定义了：
   
   // 记录着Go语言中某个数据类型的基本特征
   
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
    
在$GOROOT/src/runtime/runtime.h中定义了：

    // 有方法的interface
    struct Iface
    {
        Itab*   tab;
        void*   data;
    };
    
    // 没有方法的interface
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
    
    // interface数据类型对应的type
    type interfacetype struct {
        typ     _type
        pkgpath name
        mhdr    []imethod
    }

interface实际上是一个结构体，包括两个成员，一个是指向数据的指针，一个包含了成员的类型信息。Eface是interface{}底层使用的数据结构。因为interface中保存了类型信息，所以可以实现反射。反射其实就是查找底层数据结构的元数据。完整的实现在：$GOROOT/src/pkg/runtime/iface.c 。

## map

golang的map实现是hashtable，源码在：$GOROOT/src/runtime/hashmap.c 。

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
    
    
## unsafe.Pointer：通用指针类型，用于转换不同类型的指针，不能进行指针运算。

- unsafe.Pointer 可以和 普通指针 进行相互转换。
- unsafe.Pointer 可以和 uintptr 进行相互转换。

也就是说 unsafe.Pointer 是桥梁，可以让任意类型的指针实现相互转换，也可以将任意类型的指针转换为 uintptr 进行指针运算。
  
  
## uintptr

在golang中uintptr的定义是 type uintptr uintptr uintptr是golang的内置类型，是能存储指针的整型

用于指针运算，GC 不把 uintptr 当指针，uintptr 无法持有对象。uintptr 类型的目标会被回收。
- 一个unsafe.Pointer是一个指向变量的指针
- 但是uintptr类型的临时变量只是一个普通的数字
- 作为变量的uintptr虽然记录的是1个指针地址，但是它只是1个普通的数字，所以指向的地址很有可能被GC回收，这个无效地址空间的赋值语句将彻底摧毁整个程序！

  [1]: https://www.cnblogs.com/junneyang/p/6203710.html
