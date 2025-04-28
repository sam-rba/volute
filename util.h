#define nelem(arr) (sizeof(arr)/sizeof(arr[0]))
#define min(a, b) ((a < b) ? a : b)
#define max(a, b) ((a > b) ? a : b)

void free_arr(void **arr, int n);
int lsearch(const void *key, const void *base, size_t n, size_t size, int (*cmp)(const void *keyval, const void *datum));
