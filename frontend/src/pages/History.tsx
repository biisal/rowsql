import { useEffect, useState } from "react";
import { useSearchParams, Link } from "react-router-dom";
import {
  Clock,
  Loader2,
  ChevronLeft,
  ChevronRight,
  Database,
} from "lucide-react";
import api from "@/lib/axios";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";

interface HistoryItem {
  id: number;
  message: string;
  time: string;
}

export function History() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [history, setHistory] = useState<HistoryItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const currentPage = parseInt(searchParams.get("page") || "1");

  useEffect(() => {
    fetchHistory();
  }, [currentPage]);

  const fetchHistory = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await api.get(`/history?page=${currentPage}`);
      setHistory(response.data.data || []);
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : "Failed to fetch history";
      setError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const handlePageChange = (newPage: number) => {
    if (newPage < 1) return;
    setSearchParams({ page: newPage.toString() });
  };

  const formatTime = (timeString: string) => {
    const date = new Date(timeString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return "Just now";
    if (diffMins < 60)
      return `${diffMins} minute${diffMins > 1 ? "s" : ""} ago`;
    if (diffHours < 24)
      return `${diffHours} hour${diffHours > 1 ? "s" : ""} ago`;
    if (diffDays < 7) return `${diffDays} day${diffDays > 1 ? "s" : ""} ago`;

    return date.toLocaleDateString() + " " + date.toLocaleTimeString();
  };

  const getActionBadge = (message: string) => {
    if (message.toLowerCase().includes("created table")) {
      return (
        <Badge
          variant="default"
          className="bg-green-500/10 text-green-500 border-green-500/20"
        >
          Created
        </Badge>
      );
    }
    if (
      message.toLowerCase().includes("dropped table") ||
      message.toLowerCase().includes("deleted")
    ) {
      return (
        <Badge
          variant="default"
          className="bg-red-500/10 text-red-500 border-red-500/20"
        >
          Deleted
        </Badge>
      );
    }
    if (message.toLowerCase().includes("updated")) {
      return (
        <Badge
          variant="default"
          className="bg-blue-500/10 text-blue-500 border-blue-500/20"
        >
          Updated
        </Badge>
      );
    }
    if (message.toLowerCase().includes("inserted")) {
      return (
        <Badge
          variant="default"
          className="bg-purple-500/10 text-purple-500 border-purple-500/20"
        >
          Inserted
        </Badge>
      );
    }
    return <Badge variant="secondary">Action</Badge>;
  };

  if (loading) {
    return (
      <div className="flex-1 p-8 overflow-y-auto">
        <div className="max-w-6xl mx-auto">
          <div className="flex items-center justify-center h-64">
            <Loader2 className="w-8 h-8 animate-spin text-primary" />
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex-1 p-8 overflow-y-auto">
        <div className="max-w-6xl mx-auto">
          <Card className="border-destructive/50">
            <CardHeader>
              <CardTitle className="text-destructive">Error</CardTitle>
              <CardDescription>Failed to load history</CardDescription>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground">{error}</p>
              <Button onClick={fetchHistory} className="mt-4">
                Try Again
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 p-8 overflow-y-auto">
      <div className="max-w-6xl mx-auto space-y-6 animate-fade-in-up">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight flex items-center gap-2">
              <Clock className="w-8 h-8 text-primary" />
              History
            </h1>
            <p className="text-muted-foreground mt-1">
              View all database operations and changes
            </p>
          </div>
          <Link to="/">
            <Button variant="outline" size="sm">
              <ChevronLeft className="w-4 h-4 mr-1" />
              Back to Home
            </Button>
          </Link>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Database Operations</CardTitle>
            <CardDescription>
              Complete history of all database operations (Page {currentPage})
            </CardDescription>
          </CardHeader>
          <CardContent>
            {history.length === 0 ? (
              <div className="text-center py-12">
                <Database className="w-12 h-12 mx-auto text-muted-foreground/50 mb-4" />
                <p className="text-muted-foreground">
                  No history records found
                </p>
                <p className="text-sm text-muted-foreground/70 mt-1">
                  Database operations will appear here
                </p>
              </div>
            ) : (
              <div className="space-y-3">
                {history.map((item) => (
                  <div
                    key={item.id}
                    className="flex items-start gap-4 p-4 rounded-lg border border-border/50 hover:bg-muted/30 transition-colors"
                  >
                    <div className="flex-shrink-0 mt-1">
                      <div className="w-2 h-2 rounded-full bg-primary"></div>
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        {getActionBadge(item.message)}
                        <span className="text-xs text-muted-foreground">
                          #{item.id}
                        </span>
                      </div>
                      <p className="text-sm font-medium break-words">
                        {item.message}
                      </p>
                      <p className="text-xs text-muted-foreground mt-1 flex items-center gap-1">
                        <Clock className="w-3 h-3" />
                        {formatTime(item.time)}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        {history.length > 0 && (
          <div className="flex items-center justify-between">
            <div className="text-sm text-muted-foreground">
              Showing page {currentPage}
            </div>
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => handlePageChange(currentPage - 1)}
                disabled={currentPage === 1}
              >
                <ChevronLeft className="w-4 h-4 mr-1" />
                Previous
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => handlePageChange(currentPage + 1)}
                disabled={history.length === 0}
              >
                Next
                <ChevronRight className="w-4 h-4 ml-1" />
              </Button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
