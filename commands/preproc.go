package commands

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mitsuse/kompl/sentencizer"
	"github.com/mitsuse/kompl/tokenizer"
)

func NewPreprocCommand() cli.Command {
	command := cli.Command{
		Name:      "preproc",
		ShortName: "p",
		Usage:     "Tokenizes and sentencizes corpora.",
		Action:    preproc,

		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "ignore-blank,b",
				Usage: "Ignore blank lines.",
			},

			cli.BoolFlag{
				Name:  "ignore-xml,x",
				Usage: "Ignore xml tags.",
			},

			cli.StringFlag{
				Name:  "corpus,c",
				Value: "corpus.raw",
				Usage: "The input path of a raw corpus.",
			},

			cli.StringFlag{
				Name:  "tokenized,t",
				Value: "corpus.tokenized",
				Usage: "The output path of the tokenized corpus.",
			},
		},
	}

	return command
}

func preproc(ctx *cli.Context) {
	t := tokenizer.NewEnglishTokenizer()
	s := sentencizer.NewEnglishSentencizer()

	corpusFile, err := os.Open(ctx.String("corpus"))
	if err != nil {
		PrintError(ERROR_LOADING_CORPUS, err)
		return
	}
	defer corpusFile.Close()

	tokenizedFile, err := os.Create(ctx.String("tokenized"))
	if err != nil {
		PrintError("failed to create a tokenize corpus", err)
		return
	}
	defer tokenizedFile.Close()

	ignoreBlank := ctx.Bool("ignore-blank")
	ignoreXml := ctx.Bool("ignore-xml")

	xmlPattern := regexp.MustCompile(`^<.*>$`)

	scanner := bufio.NewScanner(corpusFile)
	for scanner.Scan() {
		line := scanner.Text()

		if ignoreBlank && len(line) == 0 {
			continue
		}

		if ignoreXml && xmlPattern.MatchString(line) {
			continue
		}

		sentenceSeq := s.Sentencize(t.Tokenize(line))
		for _, tokenSeq := range sentenceSeq {
			sentence := fmt.Sprintf("%s\n", strings.Join(tokenSeq, " "))

			if _, err := tokenizedFile.WriteString(sentence); err != nil {
				PrintError("failed to write sentence to the tokenized corpus", err)
				return
			}
		}
	}

	if err := scanner.Err(); err != nil {
		PrintError("failed to read the raw corpus", err)
		return
	}
}
