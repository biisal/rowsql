import { useEffect, useState } from 'react';
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip';
import { Link } from 'react-router-dom';
import { Button } from './ui/button';
import { Github, Star } from 'lucide-react';

async function fetchGitStar() {
	const response = await fetch('https://api.github.com/repos/biisal/rowsql');
	const data = await response.json();
	return data.stargazers_count;
}
export const GitStarButton = () => {
	const [gitStar, setGitStar] = useState(0);

	useEffect(() => {
		fetchGitStar().then(setGitStar);
	}, []);

	return (
		<Tooltip>
			<TooltipTrigger asChild>
				<Button variant="outline" asChild className="h-9 px-3">
					<Link
						to="https://github.com/biisal/rowsql"
						target="_blank"
						className="flex items-center gap-2"
					>
						<Github className="h-4 w-4" />
						<span className="text-sm font-medium">GitHub</span>
						{gitStar > 0 && (
							<>
								<div className="h-3 w-px bg-border mx-1" />{' '}
								<div className="flex items-center gap-1 text-muted-foreground">
									<Star className="h-3 w-3 fill-yellow-500" />
									<span className="mt-0.5">{gitStar}</span>
								</div>
							</>
						)}
					</Link>
				</Button>
			</TooltipTrigger>
			<TooltipContent>Start on github if you like it! </TooltipContent>
		</Tooltip>
	);
};
