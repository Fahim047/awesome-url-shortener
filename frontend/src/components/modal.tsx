import { BarChart2, X, Link2, Calendar, Clock } from "lucide-react";
import type { AnalyticsResponse } from "../../types";

interface AnalyticsModalProps {
  dialogRef: React.RefObject<HTMLDialogElement>;
  analytics: AnalyticsResponse | null;
}

export const Modal = ({ dialogRef, analytics }: AnalyticsModalProps) => {
  const closeDialog = () => {
    dialogRef?.current?.close();
  };

  return (
    <dialog
      ref={dialogRef}
      className="min-h-96 fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 rounded-2xl p-6 bg-white shadow-xl w-[90%] max-w-md [&::backdrop]:bg-black/40"
    >
      {/* Header */}
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-lg font-semibold flex items-center gap-2">
          <BarChart2 className="text-indigo-600" /> Analytics
        </h2>
        <button
          onClick={closeDialog}
          className="p-1 rounded-lg hover:bg-gray-100"
        >
          <X size={18} />
        </button>
      </div>

      {/* Body */}
      {analytics ? (
        <div className="space-y-4 text-sm text-gray-700">
          <div className="flex items-center gap-2">
            <Link2 className="text-indigo-600" size={16} />
            <span className="font-medium">Original URL:</span>
            <a
              href={analytics.long_url}
              target="_blank"
              rel="noreferrer"
              className="text-indigo-600 truncate"
            >
              {analytics.long_url}
            </a>
          </div>

          {analytics.custom_alias && (
            <p>
              <span className="font-medium">Custom Alias:</span>{" "}
              {analytics.custom_alias}
            </p>
          )}

          <p>
            <span className="font-medium">Short Key:</span>{" "}
            {analytics.short_key}
          </p>

          <p>
            <span className="font-medium">Clicks:</span> {analytics.click_count}
          </p>

          <div className="flex items-center gap-2">
            <Calendar className="text-indigo-600" size={16} />
            <span className="font-medium">Created At:</span>{" "}
            {new Date(analytics.created_at).toLocaleString()}
          </div>

          {analytics.expire_at && (
            <div className="flex items-center gap-2">
              <Clock className="text-indigo-600" size={16} />
              <span className="font-medium">Expires At:</span>{" "}
              {new Date(analytics.expire_at).toLocaleString()}
            </div>
          )}
        </div>
      ) : (
        <p className="text-gray-500">Loading analytics...</p>
      )}
    </dialog>
  );
};
