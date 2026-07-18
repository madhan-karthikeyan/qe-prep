import java.util.HashMap;
import java.util.Map;

public class BuggyCache {
    private static Map<String, byte[]> cache = new HashMap<>();

    public static byte[] compute(String key) {
        if (cache.containsKey(key)) {
            return cache.get(key);
        }
        byte[] result = new byte[1024 * 1024]; // 1 MB
        cache.put(key, result);
        return result;
    }

    public static void main(String[] args) {
        int i = 0;
        while (true) {
            compute("key-" + i);
            i++;
            if (i % 100 == 0) {
                System.out.println("Cached " + i + " entries, cache size: " + cache.size());
            }
        }
    }
}
