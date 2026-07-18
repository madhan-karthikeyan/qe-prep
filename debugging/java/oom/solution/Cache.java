import java.util.LinkedHashMap;
import java.util.Map;

public class Cache {
    private final int maxEntries;
    private final Map<String, byte[]> cache;

    public Cache(int maxEntries) {
        this.maxEntries = maxEntries;
        this.cache = new LinkedHashMap<String, byte[]>(maxEntries, 0.75f, true) {
            @Override
            protected boolean removeEldestEntry(Map.Entry<String, byte[]> eldest) {
                return size() > maxEntries;
            }
        };
    }

    public byte[] compute(String key) {
        if (cache.containsKey(key)) {
            return cache.get(key);
        }
        byte[] result = new byte[1024 * 1024];
        cache.put(key, result);
        return result;
    }

    public int size() {
        return cache.size();
    }

    public static void main(String[] args) {
        Cache cache = new Cache(100);
        int i = 0;
        while (true) {
            cache.compute("key-" + i);
            i++;
            if (i % 100 == 0) {
                System.out.println("Processed " + i + " entries, cache size: " + cache.size());
            }
        }
    }
}
