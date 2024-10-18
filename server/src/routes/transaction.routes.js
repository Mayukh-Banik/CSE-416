"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
// routes/transaction.routes.ts
const express_1 = __importDefault(require("express"));
const transaction_model_1 = __importDefault(require("../models/transaction.model"));
const router = express_1.default.Router();
// Create a new transaction
router.post('/create', (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const { sender, receiver, amount } = req.body;
        const newTransaction = new transaction_model_1.default({
            sender,
            receiver,
            amount,
        });
        yield newTransaction.save();
        res.status(201).json(newTransaction);
    }
    catch (error) {
        res.status(500).json({ message: 'Error creating transaction', error });
    }
}));
// Get all transactions
router.get('/', (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const transactions = yield transaction_model_1.default.find().populate('sender receiver', 'username email');
        res.json(transactions);
    }
    catch (error) {
        res.status(500).json({ message: 'Error retrieving transactions', error });
    }
}));
exports.default = router;
