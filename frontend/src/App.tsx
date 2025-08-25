import { useState, useRef, useEffect } from "react";
import axios from "axios";
import { Link2, Copy, BarChart2 } from "lucide-react";
import { Modal } from "./components/modal";

type TopURL = {
  short_key: string;
  click_count: number;
  expire_at?: string; // ISO date string
  created_at: string; // ISO date string
};

export default function App() {
  const [url, setUrl] = useState("");
  const [customAlias, setCustomAlias] = useState("");
  const [expireAt, setExpireAt] = useState("");
  const [shortUrl, setShortUrl] = useState<string | null>(null);
  const [shortKey, setShortKey] = useState<string | null>(null);
  const [analytics, setAnalytics] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [topUrls, setTopUrls] = useState<TopURL[]>([]);
  const dialogRef = useRef<HTMLDialogElement>(
    null
  ) as React.RefObject<HTMLDialogElement>;

  const handleShorten = async () => {
    if (!url) return;
    setLoading(true);
    try {
      const payload: Record<string, any> = { long_url: url };
      if (customAlias.trim()) payload.custom_alias = customAlias.trim();
      if (expireAt.trim()) {
        // Convert from "YYYY-MM-DDTHH:MM" to "YYYY-MM-DDTHH:MM:SSZ"
        const date = new Date(expireAt);
        payload.expire_at = date.toISOString(); // always RFC3339 UTC
      }

      const res = await axios.post(
        "http://localhost:8080/api/v1/shorten",
        payload
      );

      setShortKey(res.data.short_key);
      setShortUrl(res.data.short_url);
      setAnalytics(null);

      // reset form
      setUrl("");
      setCustomAlias("");
      setExpireAt("");
      fetchTop(); // refresh leaderboard
    } catch (err) {
      console.error(err);
      alert("Failed to shorten URL");
    } finally {
      setLoading(false);
    }
  };

  const handleCopy = () => {
    if (shortUrl) navigator.clipboard.writeText(shortUrl);
  };

  const handleViewAnalytics = async (key?: string) => {
    const targetKey = key || shortKey;
    if (!targetKey) return;
    try {
      const res = await axios.get(
        `http://localhost:8080/api/v1/analytics/${targetKey}`
      );
      setAnalytics(res.data);
      dialogRef.current?.showModal();
    } catch (err) {
      console.error(err);
      alert("Failed to fetch analytics");
    }
  };

  const fetchTop = async () => {
    try {
      const res = await axios.get<TopURL[]>("http://localhost:8080/api/v1/top");
      setTopUrls(res.data);
    } catch (err) {
      console.error("Failed to fetch top URLs", err);
    }
  };

  useEffect(() => {
    fetchTop();
  }, []);

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col items-center justify-start p-6">
      <h1 className="text-4xl md:text-5xl font-bold text-gray-800 mb-8">
        🚀 Awesome URL Shortener
      </h1>

      {/* Shortener */}
      <div className="w-full max-w-lg bg-white rounded-lg shadow p-6 space-y-6">
        {/* Long URL */}
        <div className="space-y-2">
          <label className="text-sm font-medium text-gray-700">Long URL</label>
          <div className="flex gap-2">
            <input
              type="url"
              placeholder="https://example.com/very/long/url..."
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              className="flex-1 border rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
            <button
              onClick={handleShorten}
              disabled={loading}
              className="bg-indigo-600 text-white rounded-lg px-6 py-2 font-medium hover:bg-indigo-700 disabled:opacity-50"
            >
              {loading ? "..." : "Shorten"}
            </button>
          </div>
        </div>

        {/* Advanced options */}
        <details className="rounded-lg p-4 bg-gray-50">
          <summary className="cursor-pointer font-medium text-gray-700">
            Advanced Options
          </summary>
          <div className="mt-4 space-y-4">
            {/* Custom alias */}
            <div className="space-y-2">
              <label className="text-sm font-medium text-gray-700">
                Custom Alias
              </label>
              <input
                type="text"
                placeholder="e.g. my-custom-link"
                value={customAlias}
                onChange={(e) => setCustomAlias(e.target.value)}
                className="w-full border rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
              />
            </div>

            {/* Expiration */}
            <div className="space-y-2">
              <label className="text-sm font-medium text-gray-700">
                Expire At
              </label>
              <input
                type="datetime-local"
                value={expireAt}
                onChange={(e) => setExpireAt(e.target.value)}
                className="w-full border rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500"
              />
            </div>
          </div>
        </details>

        {/* Shortened URL preview */}
        {shortUrl && (
          <div className="rounded-lg border border-indigo-200 bg-indigo-50 p-4 flex items-center justify-between">
            <a
              href={shortUrl}
              target="_blank"
              rel="noreferrer"
              className="text-indigo-700 font-medium flex items-center gap-2"
            >
              <Link2 size={18} /> {shortUrl}
            </a>
            <div className="flex gap-2">
              <button
                onClick={handleCopy}
                className="p-2 hover:bg-indigo-100 rounded-lg"
                title="Copy"
              >
                <Copy size={18} />
              </button>
              <button
                onClick={() => handleViewAnalytics()}
                className="p-2 hover:bg-indigo-100 rounded-lg"
                title="View Analytics"
              >
                <BarChart2 size={18} />
              </button>
            </div>
          </div>
        )}
      </div>

      {/* Leaderboard */}
      <div className="mt-6 w-full max-w-lg bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
          <BarChart2 size={20} /> Leaderboard
        </h2>
        {topUrls.length === 0 ? (
          <p className="text-gray-500">No data yet</p>
        ) : (
          <ul className="space-y-2">
            {topUrls.map((u) => (
              <li
                key={u.short_key}
                className="flex items-center justify-between bg-gray-50 px-3 py-2 rounded-lg hover:bg-gray-100"
              >
                <a
                  href={`http://localhost:8080/${u.short_key}`}
                  target="_blank"
                  rel="noreferrer"
                  className="flex items-center gap-2 text-indigo-600 font-medium truncate"
                >
                  <Link2 size={16} />/{u.short_key}
                </a>
                <div className="flex items-center gap-3">
                  <span className="text-sm text-gray-700 font-semibold">
                    {u.click_count}
                  </span>
                  <button
                    onClick={() => handleViewAnalytics(u.short_key)}
                    className="p-1 hover:bg-indigo-100 rounded-lg"
                    title="View Analytics"
                  >
                    <BarChart2 size={18} />
                  </button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>

      <Modal dialogRef={dialogRef} analytics={analytics} />
    </div>
  );
}
