import { Link } from "react-router-dom";
import { Home, ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export function NotFound() {
  return (
    <div className="flex items-center justify-center min-h-[calc(100vh-8rem)] p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div className="mx-auto mb-4 text-6xl font-bold text-primary">
            404
          </div>
          <CardTitle className="text-2xl">Page Not Found</CardTitle>
          <CardDescription>
            The page you're looking for doesn't exist or has been moved.
          </CardDescription>
        </CardHeader>
        <CardContent className="text-center text-muted-foreground">
          <p>
            The URL you entered may be incorrect, or the resource may have been
            deleted.
          </p>
        </CardContent>
        <CardFooter className="flex flex-col sm:flex-row gap-2 justify-center">
          <Button asChild variant="default">
            <Link to="/">
              <Home className="mr-2 h-4 w-4" />
              Go to Home
            </Link>
          </Button>
          <Button asChild variant="outline" onClick={() => window.history.back()}>
            <span className="cursor-pointer">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Go Back
            </span>
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}
