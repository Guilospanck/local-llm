import React, { useEffect, useRef, useState } from "react";
import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import "./App.css";

const Search = ({
  thinking,
  onInputChange,
  onAbortClick,
  onSearchClick,
  onKeyDown,
  inputRef,
}: {
  thinking: boolean;
  onInputChange: (val: string) => void;
  onAbortClick: () => void;
  onSearchClick: () => void;
  onKeyDown: (e: React.KeyboardEvent<HTMLInputElement>) => void;
  inputRef: React.Ref<HTMLInputElement>;
}) => {
  return (
    <div className="search-container">
      <div className="info">{thinking ? "Thinking..." : ""}</div>
      <div className="search">
        <input
          ref={inputRef}
          className="input"
          disabled={thinking}
          type="search"
          onKeyDown={onKeyDown}
          onChange={(e) => onInputChange(e.target.value)}
        />
        <button onClick={onSearchClick}>Search</button>
        <button disabled={!thinking} onClick={onAbortClick}>
          Abort
        </button>
      </div>
    </div>
  );
};

function App() {
  const [query, setQuery] = useState<string>("");
  const [thinking, setThinking] = useState<boolean>(false);
  const [messageElements, setMessageElements] = useState<
    Array<{ answer: React.ReactNode; time?: string; question: React.ReactNode }>
  >([]);
  const [controller, setController] = useState<AbortController | null>(null);
  const inputRef = useRef<HTMLInputElement | null>(null);
  const [currentAnswer, setCurrentAnswer] = useState<string>("");

  const focusInput = () => {
    if (inputRef.current) {
      inputRef.current.focus();
    }
  };

  useEffect(() => {
    focusInput();
  }, []);

  const createMarkdownElement = (data: string) => {
    return (
      <Markdown remarkPlugins={[remarkGfm]} className="message">
        {data}
      </Markdown>
    );
  };

  const onKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.code === "Enter") {
      onSearchClick();
    }
  };

  const removeThinkTag = (data: string) =>
    data.replace("<think>", "").replace("</think>", "");

  const onSearchClick = async () => {
    const t0 = performance.now();

    setThinking(true);

    const controller = new AbortController();
    setController(controller);

    let data = "";

    const response = await fetch("http://localhost:4444/", {
      signal: controller.signal,
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ query }),
    });

    if (!response.body) {
      console.error("ReadableStream not supported");
      return;
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder();

    while (true) {
      const { value, done } = await reader.read();
      if (done) {
        break;
      }

      const chunk = decoder.decode(value, { stream: true });
      if (!chunk) continue;

      data += chunk;

      setCurrentAnswer(removeThinkTag(data));
    }

    data = removeThinkTag(data);

    // calculate benchmark
    const t1 = performance.now();
    const seconds = (t1 - t0) / 1_000;
    const secondsFormatted = `Took ${seconds.toFixed(2)} seconds`;

    setMessageElements((prev) => [
      ...prev,
      {
        answer: createMarkdownElement(data),
        time: secondsFormatted,
        question: createMarkdownElement(query),
      },
    ]);

    setThinking(false);
    setCurrentAnswer("");
  };

  const onAbortClick = () => {
    controller?.abort();
    setThinking(false);
  };

  const onInputChange = (val: string) => {
    setQuery(val);
  };

  return (
    <>
      <div className="container">
        <div id="messages" className="messages">
          {messageElements.map((elm, index) => (
            <React.Fragment key={index}>
              <div className="right-container">
                <div className="right">{elm.question}</div>
              </div>
              <div className="left">{elm.answer}</div>
              {elm.time && <div className="info">{elm.time}</div>}
            </React.Fragment>
          ))}

          {thinking && (
            <React.Fragment>
              <div className="right-container">
                <div className="right">{createMarkdownElement(query)}</div>
              </div>
              <div className="left">{createMarkdownElement(currentAnswer)}</div>
            </React.Fragment>
          )}
        </div>
        <Search
          inputRef={inputRef}
          thinking={thinking}
          onSearchClick={onSearchClick}
          onAbortClick={onAbortClick}
          onInputChange={onInputChange}
          onKeyDown={onKeyDown}
        />
      </div>
    </>
  );
}

export default App;
