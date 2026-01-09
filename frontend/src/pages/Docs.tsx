import { BookOpen, ArrowLeft } from 'lucide-react';
import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';

export function Docs() {
	return (
		<div className="flex-1 p-8 overflow-y-auto">
			<div className="max-w-2xl mx-auto">
				<div className="flex flex-col items-center justify-center min-h-[60vh] space-y-6 text-center">
					<BookOpen className="w-16 h-16 text-muted-foreground/50" />
					<div className="space-y-2">
						<h1 className="text-3xl font-bold">Documentation</h1>
						<p className="text-lg text-muted-foreground">
							Still writing... Updating soon
						</p>
					</div>
					<Link to="/">
						<Button variant="outline" className="gap-2">
							<ArrowLeft className="w-4 h-4" />
							Back to Home
						</Button>
					</Link>
				</div>
			</div>
		</div>
	);
}
