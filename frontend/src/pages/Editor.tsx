import { useRef, useState, useEffect } from "react";

import type { ErrorModel, RunSqlQueryOutputBody } from "@/client";
import { listTablesOptions, runSqlQueryMutation } from "@/client/@tanstack/react-query.gen";
import { Editor as SqlEditor, type Monaco, type OnMount } from "@monaco-editor/react";
import type * as monacoEditor from 'monaco-editor';
import { useMutation, useQuery } from "@tanstack/react-query";
import {
    ResizableHandle,
    ResizablePanel,
    ResizablePanelGroup,
} from "@/components/ui/resizable";
import { EditorHeader } from "@/components/editor/EditorHeader";
import { QueryTable } from "@/components/editor/QueryTable";


export const Editor = () => {
    const editorRef = useRef<Monaco>(null);
    const [queryResult, setQueryResult] = useState<RunSqlQueryOutputBody>();
    const [queryError, setQueryError] = useState<string>("");
    const [runText, setRunText] = useState<"Run" | "Run Selected">("Run");
    const { data: tables } = useQuery(listTablesOptions());
    const { refetch } = useQuery(listTablesOptions());

    const tablesRef = useRef(tables);
    useEffect(() => {
        tablesRef.current = tables;
    }, [tables]);

    const editorQueryMutation = useMutation(runSqlQueryMutation());

    const handleEditorDidMount: OnMount = (editor, monaco) => {
        editorRef.current = editor;
        editor.onDidChangeCursorSelection((e) => {
            setRunText(e.selection.isEmpty() ? "Run" : "Run Selected");
        });

        editor.addCommand(
            monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter,
            () => handleRunQuery()
        );
        monaco.languages.registerCompletionItemProvider("sql", {
            triggerCharacters: [" ", "."],
            provideCompletionItems: (
                model: monacoEditor.editor.ITextModel,
                position: monacoEditor.Position
            ): monacoEditor.languages.ProviderResult<monacoEditor.languages.CompletionList> => {
                const word = model.getWordUntilPosition(position);
                const range: monacoEditor.IRange = {
                    startLineNumber: position.lineNumber,
                    endLineNumber: position.lineNumber,
                    startColumn: word.startColumn,
                    endColumn: word.endColumn,
                };

                const currentTables = tablesRef.current;
                if (!currentTables) return { suggestions: [] };

                const suggestions: monacoEditor.languages.CompletionItem[] = [
                    ...currentTables.map((table) => ({
                        label: table.tableName,
                        kind: monaco.languages.CompletionItemKind.Class,
                        insertText: table.tableName,
                        range,
                    })),
                ];

                return { suggestions };
            },
        });
    };

    const getQueryText = (): string | undefined => {
        const editor = editorRef.current as Monaco | null;
        if (!editor) return;
        const model = editor.getModel();
        const selection = editor.getSelection();
        return model.getValueInRange(selection) || model.getValue();
    };

    const handleRunQuery = () => {
        const query = getQueryText();
        if (!query) return;

        editorQueryMutation.mutate(
            { body: { query } },
            {
                onSuccess: (data) => {
                    setQueryError("");
                    setQueryResult(data);
                    if (query.toLowerCase().includes("table")) {
                        refetch();
                    }
                },
                onError: (error) => {
                    const err = error as unknown as ErrorModel;
                    setQueryError(
                        err.errors?.[0]?.message || err.detail || "An unknown error occurred"
                    );
                },
            }
        );
    };

    return (
        <div className="h-full w-full bg-background font-mono text-sm">
            <ResizablePanelGroup orientation="vertical">
                <ResizablePanel defaultSize="55%">
                    <div className="flex h-full flex-col">
                        <EditorHeader
                            runText={runText}
                            loading={editorQueryMutation.isPending}
                            onRun={handleRunQuery}
                        />
                        <div className="flex-1">
                            <SqlEditor
                                height="100%"
                                theme="vs-dark"
                                language="sql"
                                defaultLanguage="sql"
                                onMount={handleEditorDidMount}
                                defaultValue="-- Start typing your SQL query here..."
                                options={{
                                    fontSize: 14,
                                    lineHeight: 22,
                                    fontFamily: "JetBrains Mono, Fira Code, monospace",
                                    minimap: { enabled: false },
                                    scrollBeyondLastLine: false,
                                    renderLineHighlight: "gutter",
                                    padding: { top: 12, bottom: 12 },
                                    scrollbar: { vertical: "auto", horizontal: "auto" },
                                }}
                            />
                        </div>
                    </div>
                </ResizablePanel>

                <ResizableHandle withHandle />

                <ResizablePanel>
                    <QueryTable data={queryResult} error={queryError} />
                </ResizablePanel>
            </ResizablePanelGroup>
        </div>
    );
};


